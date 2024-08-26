package go_license_checker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type LicenseChecker struct {
	AllowedSNs      []string
	AllowedSNPrefix string
}

func init() {
	allowedSNs := []string{
		"b2a9be8c7159d8de",
		"f28fa9cb89609e6e",
		// 在这里添加更多允许的SN
	}
	// 将允许的SN列表中的所有SN转换为小写
	for i, sn := range allowedSNs {
		allowedSNs[i] = strings.ToLower(sn)
	}
	licenseChecker = &LicenseChecker{
		AllowedSNs:      allowedSNs,
		AllowedSNPrefix: "lincosdemo",
	}
}

var licenseChecker *LicenseChecker

func Check() error {
	go func() {
		ticker := time.NewTicker(time.Hour) // 每小时检查一次
		defer ticker.Stop()

		for {
			err := licenseChecker.checkLicense()
			if err != nil {
				fmt.Println("License check failed:", err)
				logAndExit("License check failed: " + err.Error())
			}
			<-ticker.C // 等待下一个检查时间
		}
	}()
	return nil
}

func (lc *LicenseChecker) checkLicense() error {
	// 获取设备的SN
	sn, err := getSerialNumber()
	if err != nil {
		return err
	}

	// 检查SN是否在允许的列表中或以允许的前缀开头
	if !lc.isAllowedSN(sn) {
		return fmt.Errorf("device SN %s is not allowed", sn)
	}

	// 检查当前时间是否在2024年内
	now := time.Now()
	if now.Year() != 2024 {
		return fmt.Errorf("license is only valid in 2024")
	}

	return nil
}

func (lc *LicenseChecker) isAllowedSN(sn string) bool {
	// 检查SN是否在允许的列表中
	for _, allowedSN := range lc.AllowedSNs {
		if sn == allowedSN {
			return true
		}
	}

	// 检查SN是否以允许的前缀开头
	return strings.HasPrefix(strings.ToLower(sn), strings.ToLower(lc.AllowedSNPrefix))
}

func getSerialNumber() (string, error) {
	// 首先尝试使用getprop获取SN
	sn, err := getSerialNumberFromGetprop()
	if err == nil && sn != "" {
		return strings.ToLower(sn), nil
	}

	// 如果getprop失败或返回空字符串，尝试使用DMI信息获取SN
	sn, err = getSerialNumberFromDMI()
	if err == nil && sn != "" {
		return strings.ToLower(sn), nil
	}

	// 如果两个方法都失败，返回错误
	return "", fmt.Errorf("failed to get serial number")
}

func getSerialNumberFromGetprop() (string, error) {
	out, err := exec.Command("getprop", "ro.serialno").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getSerialNumberFromDMI() (string, error) {
	// DMI信息文件的路径
	dmiPath := "/sys/class/dmi/id/product_serial"

	// 读取文件内容
	content, err := os.ReadFile(dmiPath)
	if err != nil {
		// 如果出现错误，返回错误信息
		return "", err
	}

	// 如果成功，转换内容为字符串并去除空格和换行
	serialNumber := strings.TrimSpace(string(content))

	// 返回序列号和 nil 表示没有错误
	return serialNumber, nil
}

func logAndExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

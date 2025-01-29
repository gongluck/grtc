/*
 * @Author: gongluck
 * @Date: 2025-01-29 20:05:19
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-29 20:06:11
 */

package util

import (
	"log"
	"net"
)

// 获取唯一设备 ID（MAC 地址）
func GetDeviceID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error getting network interfaces: %v", err)
	}

	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String() // 返回 MAC 地址
		}
	}

	log.Fatal("No valid network interface found.")
	return ""
}

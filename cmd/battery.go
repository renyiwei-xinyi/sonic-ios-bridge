/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/mitchellh/mapstructure"
	"os"

	"github.com/spf13/cobra"
)

var batteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Show battery of your device.",
	Long:  "Show battery of your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(udid) != 0 {
			device := util.GetDeviceByUdId(udid)
			if device == nil {
				os.Exit(0)
			}
			b := entity.Battery{}
			bd, err := device.GetValue("com.apple.mobile.battery", "")
			if err != nil {
				return util.NewErrorPrint(util.ErrSendCommand, "get value", err)
			}
			bi := entity.BatteryInter{}
			mapstructure.Decode(bd, &bi)
			b.SerialNumber = device.Properties().SerialNumber
			b.Level = bi.BatteryCurrentCapacity
			b.Temperature = 0
			data := util.ResultData(b)
			fmt.Println(util.Format(data, isFormat, isJson))
		} else {
			usbMuxClient, err := giDevice.NewUsbmux()
			if err != nil {
				return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
			}
			list, err1 := usbMuxClient.Devices()
			if err1 != nil {
				return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
			}
			if len(list) != 0 {
				var batteryList entity.BatteryList
				for _, d := range list {
					b := entity.Battery{}
					bd, err := d.GetValue("com.apple.mobile.battery", "")
					if err != nil {
						continue
					}
					bi := entity.BatteryInter{}
					mapstructure.Decode(bd, &bi)
					b.SerialNumber = d.Properties().SerialNumber
					b.Level = bi.BatteryCurrentCapacity
					b.Temperature = 0
					batteryList.BatteryInfo = append(batteryList.BatteryInfo, b)
				}
				data := util.ResultData(batteryList)
				fmt.Println(util.Format(data, isFormat, isJson))
			} else {
				fmt.Println("no device connected")
				os.Exit(0)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(batteryCmd)
	batteryCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	batteryCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	batteryCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}

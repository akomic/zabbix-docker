package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"zabbix-docker/cadvisor"
)

var (
	config  string
	daemon  bool
	version bool

	// Cmd ...
	Cmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		Run:   run,
	}
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to figure out our hostname")
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig)

	Cmd.PersistentFlags().StringP("config", "c", "", "config file (default: ~/.zabbix-docker/config.yml)")
	Cmd.PersistentFlags().StringP("addr", "", "127.0.0.1:8080", "CAdvisor server address")
	Cmd.PersistentFlags().StringP("zabbixAddr", "", "127.0.0.1:10051", "Zabbix Server address")
	Cmd.PersistentFlags().StringP("hostname", "", hostname, "Hostname (for Zabbix)")
	Cmd.PersistentFlags().StringP("hostGroup1", "", "Docker Containers", "Docker Label for Prototype Hostgroup")
	Cmd.PersistentFlags().StringP("hostGroup2", "", "Docker Containers", "Docker Label for Prototype Hostgroup")
	Cmd.PersistentFlags().StringP("hostGroup3", "", "Docker Containers", "Docker Label for Prototype Hostgroup")
	Cmd.PersistentFlags().StringP("hostGroup4", "", "Docker Containers", "Docker Label for Prototype Hostgroup")
	Cmd.PersistentFlags().BoolP("zabbix", "z", false, "Discover and send to Zabbix")
	Cmd.PersistentFlags().BoolP("discovery", "d", false, "Discovery")
	Cmd.PersistentFlags().StringP("containerId", "i", "", "Container ID")
	Cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbosity")

	Cmd.PersistentFlags().MarkHidden("addr")
	Cmd.PersistentFlags().MarkHidden("zabbixAddr")
	Cmd.PersistentFlags().MarkHidden("hostname")
	Cmd.PersistentFlags().MarkHidden("hostGroup1")
	Cmd.PersistentFlags().MarkHidden("hostGroup2")
	Cmd.PersistentFlags().MarkHidden("hostGroup3")
	Cmd.PersistentFlags().MarkHidden("hostGroup4")

	viper.BindPFlag("config", Cmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("addr", Cmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("zabbixAddr", Cmd.PersistentFlags().Lookup("zabbixAddr"))
	viper.BindPFlag("hostname", Cmd.PersistentFlags().Lookup("hostname"))
	viper.BindPFlag("hostGroup1", Cmd.PersistentFlags().Lookup("hostGroup1"))
	viper.BindPFlag("hostGroup2", Cmd.PersistentFlags().Lookup("hostGroup2"))
	viper.BindPFlag("hostGroup3", Cmd.PersistentFlags().Lookup("hostGroup3"))
	viper.BindPFlag("hostGroup4", Cmd.PersistentFlags().Lookup("hostGroup4"))
	viper.BindPFlag("zabbix", Cmd.PersistentFlags().Lookup("zabbix"))
	viper.BindPFlag("discovery", Cmd.PersistentFlags().Lookup("discovery"))
	viper.BindPFlag("containerId", Cmd.PersistentFlags().Lookup("containerId"))
	viper.BindPFlag("verbose", Cmd.PersistentFlags().Lookup("verbose"))

	viper.BindEnv("addr", "CADVISOR_HTTP_ADDR")
	viper.BindEnv("addr", "ZABBIX_ADDR")
}

func initConfig() {
	cfgFile := viper.GetString("config")

	if cfgFile == "" {
		abs, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), ".zabbix-docker/config.yml"))
		if err == nil {
			cfgFile = abs
		}
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read config file: ", err.Error())
		}
	}

	addr := viper.GetString("addr")

	if addr == "" {
		fmt.Println("You need to configure access to CAdvisor through: config file/env/flags")
		os.Exit(1)
	}
}

func run(ccmd *cobra.Command, args []string) {
	zabbixAddr := viper.GetString("zabbixAddr")

	zabbix := viper.GetBool("zabbix")
	discovery := viper.GetBool("discovery")
	containerId := viper.GetString("containerId")
	switch {
	case zabbix:
		if zabbixAddr == "" {
			fmt.Println("You need to configure access to Zabbix through: config file/env/flags")
			os.Exit(1)
		}
		cadvisor.Zabbix()
	case discovery:
		cadvisor.Containers()
	case containerId != "":
		cadvisor.Container()
	default:
		ccmd.HelpFunc()(ccmd, args)
	}
}

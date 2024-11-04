package nosql

import (
	"fmt"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func (s *Storage) SaveLocale(locale string) error {
	cfg, err := s.GetGuiConfig()
	if err != nil {
		return fmt.Errorf(`read config: %w`, err)
	}
	cfg.Locale = locale
	err = s.writeConfig(cfg)
	if err != nil {
		return fmt.Errorf(`save locale: %w`, err)
	}
	return nil
}

func (s *Storage) SaveTheme(theme model.ThemeName) error {
	cfg, err := s.GetGuiConfig()
	if err != nil {
		return fmt.Errorf(`read config: %w`, err)
	}
	cfg.Theme = theme
	err = s.writeConfig(cfg)
	if err != nil {
		return fmt.Errorf(`save locale: %w`, err)
	}
	return nil
}

func (s *Storage) SaveSidebarMinimized(sidebarMinified bool) error {
	cfg, err := s.GetGuiConfig()
	if err != nil {
		return fmt.Errorf(`read config: %w`, err)
	}
	cfg.SideBarMinimized = sidebarMinified
	err = s.writeConfig(cfg)
	if err != nil {
		return fmt.Errorf(`save sidebarMinified: %w`, err)
	}
	return nil
}

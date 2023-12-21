package nosql

import "fmt"

func (s *Storage) SaveLocale(locale string) error {
	cfg, err := s.GetGuiConfig()
	if err != nil {
		return fmt.Errorf(`failed to read config: %w`, err)
	}
	cfg.Locale = locale
	err = s.writeConfig(cfg)
	if err != nil {
		return fmt.Errorf(`failed to save locale: %w`, err)
	}
	return nil
}

func (s *Storage) SaveSidebarMinimized(sidebarMinified bool) error {
	cfg, err := s.GetGuiConfig()
	if err != nil {
		return fmt.Errorf(`failed to read config: %w`, err)
	}
	cfg.SideBarMinimized = sidebarMinified
	err = s.writeConfig(cfg)
	if err != nil {
		return fmt.Errorf(`failed to save sidebarMinified: %w`, err)
	}
	return nil
}

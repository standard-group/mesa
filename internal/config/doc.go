// Package config centralizes all Mesa server configuration.
// Values come from (in order of precedence):
//  1. Environment variables
//  2. TOML file pointed to by MESA_CONFIG_PATH
//  3. Built-in defaults
package config

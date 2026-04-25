// Package config handles loading and validating logpipe configuration files.
//
// Configuration is expressed as YAML and supports the following top-level keys:
//
//	# sources lists the log files logpipe should tail.
//	sources:
//	  - name: myapp          # optional human-readable label
//	    path: /var/log/app.log
//
//	# min_level drops log lines whose "level" field is below this value.
//	# Recognised levels (ascending): trace, debug, info, warn, error, fatal.
//	min_level: info
//
//	# format controls output rendering: "raw" (default) or "pretty".
//	format: pretty
//
//	# key_filter retains only lines where the given key is present in the JSON.
//	key_filter: request_id
//
// Use [Load] to read a config file from disk.
package config

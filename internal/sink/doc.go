// Package sink provides the final stage of the logpipe processing pipeline.
//
// A Sink consumes formatted log lines delivered over a string channel and
// writes them to an io.Writer — typically os.Stdout or an open log file.
//
// Basic usage:
//
//	s := sink.New(os.Stdout)
//	if err := s.Run(ctx, formattedLines); err != nil {
//		log.Fatal(err)
//	}
//
// To write to a file instead of stdout, use FileWriter to obtain a writer
// and a close function:
//
//	w, closeFn, err := sink.FileWriter("/var/log/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer closeFn()
//	s := sink.New(w)

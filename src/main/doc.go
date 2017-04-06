package main

/*

  +----------+
  |          |-+
  | internal | |-+   Common util packages
  |          | | |
  +----------+ | |
   +-----------+ |
     +-----------+

  +----------+
  |          |-+
  |   main   | |-+  Application domain packages
  |          | | |
  +----------+ | |
   +-----------+ |
     +-----------+

  +----------+
  |          |-+
  |  vendor  | |-+  Third-party packages
  |          | | |
  +----------+ | |
   +-----------+ |
     +-----------+
*/

// common pattern:
// main -> command(flags, env) -> some deps -> router -> api -> server.ListenAndServe()

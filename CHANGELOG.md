# Changelog

All notable changes to this project will be documented in this file.

## `v0.3.1B` - 2025-10-28

Fixed some memory handling bugs in basic rituals.
Created new spell: `x86linux/evasion/reverse_shell`: the same thing than basic/reverse_shell but with evasion techniques to avoid detection by antivirus, avoid execution in sandboxes and avoid debugging.

## `v0.3.0B` - 2025-10-22

-   **Features:**
    *   Introduced a new Malfoy chamber for runtime spell generation with custom parameters.
    *   Enhanced connection and command handling, including fixes for multiple connections, improved bot data formatting, and robust argument parsing.
    *   Added spell listing functionality and refactored for new "marauder" integration.
    *   Implemented debug printing and `SPrintScroll` function.

-   **Bug Fixes:**
    *   Corrected argument order for rune listing.
    *   Fixed incorrect base/size issues in `ParseInt`.
    *   Resolved several breaking issues in `main.c`.
    *   Addressed a typo in the build description.

-   **Refactorings:**
    *   Optimized TCP ritual for efficiency and compatibility with new "marauder/imperius" components.
    *   Updated "imperius" and "marauder" for better compatibility and new connection handling.
    *   Revised Ritual structure to support Listener and Connect types, improving overall functionality.
    *   Consolidated wrapper libraries (e.g., `-lfidelius/rituals/utils` to `-lwrapper`).

-   **Documentation:**
    *   Updated TODO list.
    *   Improved clarity of existing comments.

-   **Chores:**
    *   Upgraded Go to version 1.25.1.
    *   Added `_tmp` to `.gitignore`.

-   **Style:**
    *   Applied consistent indentation in `copy.go`.

## `v0.2.0B` - 2025-08-25

- Renamed the project from "onTop" to "The Dark Mark".
- Refactored codebase for theme consistency.
- Created a new banner and updated all relevant documentation.
- Created Fidelius, this is a new feature that is a encoder/decoder for the communication between the client and the server, it can be used in payload encoding and decoding too.
- Created Rituals, this is a new feature that is the possibility to create custom protocols for the communication between the client and the server.

## `v0.1.2B` - 2025-08-20

- Increased interactivity of the pseudo shell.
- Revamped README usage section for new commands.

## `v0.1.1B` - 2025-08-19

- Improved maintenance and management of modules.
- Enhanced `session` command to be able to leave the session.
- Added `use` command to select modules for execution.
- Added `options` command to view and manage module options.
- Added `set` command to modify module options.
- Refactored frontend to improve user experience.
- Fixed various bugs and improved stability.

## `v0.1.0B` - 2025-08-02

- Initial release of the onTop C2 framework.
- Basic functionality for client-server communication.
- Support for executing commands on connected clients.
- Added basic session management.
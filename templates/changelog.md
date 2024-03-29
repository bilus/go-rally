# Changelog

## v0.5.1

### Fixed

- Fix 500 instead of validation error message when saving a post with empty
  body.
- Fix panic while trying to delete a board.


## v0.5.0

### Added

- Show comment count next to posts listed for a board.


## v0.4.1

### Fixed

- Fix comments cannot be deleted.


## v0.4.0

### Added

- Post reactions can now be added/removed from the board itself.
- Posts can be archived (permanent deletion is still an option); users have an
  option to see all archived posts for a board.

### Fixed

- Cosmetic visual glitches.


## v0.3.1

### Fixed

- Fix dashboard crashing for users without own boards.


## v0.3.0

### Added

- Slack-style reactions to posts.
- Support for starring boards.
- Boards can be made private, shared using direct links.
- User dashboard showing the user's own and starred boards.
- A drop-down for board selection to the navigation bar.
- "Insert emoji" button to markdown editor toolbar.
- Markdown editor for comments.
- Markdown code blocks are rendered with syntax highlighting.
- A number of minor improvements, including ones to readability, look & feel.


## v0.2.1

### Fixed

- Fix discrepancy of number of votes per post.


## v0.2.0

### Added

- Logged in users can create any number of boards.
- Easy switching between boards; remember the last visited board.
- Board owners can manage their boards.
- Per-board vote budget, board owners can refill votes.
- Better layout.


## v0.1.0

### Added

- Login using email and password or a Google account.
- One global board.
- Logged in users can create posts and comments.
- Authors can manage their posts and comments.
- Basic upvoting and downvoting.
- Sorting posts: Top/newest.
- Markdown post editor with support for pasting images.
- Support for post drafts.
- Anonymous posts and comments.

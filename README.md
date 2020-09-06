# Welcome to Buffalo!

Thank you for choosing Buffalo for your web development needs.

## Database Setup

It looks like you chose to set up your application using a database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start it for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started your database, now Buffalo can create the databases in that file for you:

	$ buffalo pop create -a

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ oya run buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page. You can run the following command to open it:

    $ oya run browse

**Congratulations!** You now have your Buffalo application up and running.

## What Next?

We recommend you heading over to [http://gobuffalo.io](http://gobuffalo.io) and reviewing all of the great documentation there.

Good luck!

[Powered by Buffalo](http://gobuffalo.io)

## TODO

- [X] Login using form
- [X] Login using Google actually logs you in
- [x] Posts pages are protected
- [X] Logout
- [X] Make tests pass
- [X] Create new users when signing on using Google
- [X] Option to restrict Google sign in to a particular domain
- [X] Store Google UserID and use it instead of email to match
- [X] Option to disable sign ups using env
- [X] Only author can delete/edit their message
- [X] Hide controls for unauthorized users
- [X] Downvoting
- [X] Downvoting only if upvoted
- [X] Upvoting up to the available budget
- [X] Show user's available votes
- [X] Top/new
- [X] Show time ago post created
- [X] Good embedded editor
- [X] Support pasting images
    - [X] PoC
    - [X] Store images in a configurable dir
    - [X] Links to images -> show
    - [X] Delete attachment model -> delete file
    - [X] BLOCKED: Drafts implemented
    - [X] Tests
- [X] Drafts
    - [X] New posts are immediately created as drafts
    - [X] Nobody sees drafts
    - [X] Draft can be saved
    - [X] Draft can be published, post can be unpublished
    - [X] Draft can be deleted while editing
    - [X] Author can see a list of their draft posts
    - [X] Author can edit/delete drafts
- [X] Anonymous posts
- [X] Comments
    - [X] List under post
    - [X] Adding -> show
    - [X] Show author email for listed comments
    - [X] Basic markdown (no editor)
    - [X] Author can delete
    - [X] Anonymous
    - [X] Spinner when loading
- [X] Styling
- [X] Deploy
    - [X] Create heroku app
    - [X] Set up env variables
    - [X] OAuth callback
    - [X] Attachments in PG
- [ ] Boards
    - [X] Board resource
    - [X] User can create a new board
    - [X] Posts are associated with a board
    - [X] Switching between boards
    - [X] Remove drafts
    - [X] Remove unnecessary routes, handlers
      Leave listing posts for administrative purposes but remove route for now.
    - [X] Remember last dashboard
    - [X] Vote budget is per-board, not global
    - [X] Clean up code
      - [X] Update post votes based on count from upvote/downvote in /votes create/destroy.
      - [X] Failing tests
      - [X] Remove User#Votes.
      - [X] Real audit.
      - [X] Remove Vote.
    - [X] When creating/editing board, you specify voting policy
       - [X] Board limit
    - [X] Allow empty max votes
    - [X] Board creator can edit/destroy board
    - [ ] Board creator can reset votes
    - [ ] Style everything
      - [ ] Try to use morphdom to fix flicker when upvoting/downvoting
    - [ ] Test everything
- [ ] Prepare for Tooploox test run
    - [ ] Redis configuration from open redis var
    - [ ] Upvote post on its page
    - [ ] App version
    - [ ] Deploy
    - [ ] Test
- [ ] Voting improvements
    - [ ] Can take votes back for X minutes
    - [ ] Vote limit can be turned off
    - [ ] Vote reefill schedule.
- [ ] Board management
    - [ ] Board owner can edit/delete the board (= creator)
    - [ ] Board records are only marked as deleted
    - [ ] Owners can add other owners
    - [ ] Board list
    - [ ] Changing board title
    - [ ] Posting rights can be limited to owners
- [ ] Moderation
- [ ] Scheduled job to reset available votes
    - [ ] PoC job
    - [ ] Configurable schedule and num of votes
    - [X] New users have votes
- [ ] Projects/tracking progress
    - [ ] User defined pipelines?
      Board A -> Board B -> Board C
      Pipeline view
- [ ] Private boards
- [ ] Reactions
      See a list of who reacted, slack-like
      Users can define conventions e.g. "I want to help"
- [ ] Tenancy
- [ ] Use websockets to show user activity in real-time
      "Gina just upvoted a post"
- [ ] User can see their action history
      Admin can see it too. (Anonymous actions?)
- [ ] Admin can delete/edit any posts
- [ ] BUG: Navbar collapsing but no hamburger menu
- [ ] Slack/email integration
  beehive? https://github.com/muesli/beehive
- [ ] Better urls
    - [ ] Board slugs
    - [ ] Post slugs

## Deployment

OAuth callback: https://console.developers.google.com/apis/credentials?pli=1

## Template

1. https://demos.creative-tim.com/material-kit/docs/2.1/getting-started/introduction.html?_ga=2.91541151.440348291.1598815666-986970519.159881566
6

## Font Awesome icons

A list of icons we use:

- `fa-caret-down`
- `fa-caret-up`
- `fa-google`
- `fa-home`
- `fa-lightbulb`
- `fa-plus`
- `fa-question-circle`
- `fa-trash-alt`

To minimize the size of the CSS, we build a custom one using https://icomoon.io/app/. If you want to add another icon, rebuild the `font-awesome.css` (see [The Easy Way](https://blog.webjeda.com/optimize-fontawesome/)).

## Design ideas

1. https://www.vectorstock.com/royalty-free-vector/brainstorm-and-creative-idea-concept-vector-20635223
1. https://www.vectorstock.com/royalty-free-vectors/bulb-puzzle-vectors

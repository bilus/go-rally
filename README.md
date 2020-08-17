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
- [ ] Top/new
- [ ] Good embedded editor
- [ ] Tenancy
- [ ] Styling
- [ ] Admin can delete/edit any posts
- [ ] Scheduled job to reset available votes

## Design ideas

1. http://preview.viui18.com/clarity/index.html
1. https://colorlib.com/preview/theme/safario/
1. https://colorlib.com/preview/#podcast
1. https://www.vectorstock.com/royalty-free-vector/brainstorm-and-creative-idea-concept-vector-20635223
1. https://www.vectorstock.com/royalty-free-vectors/bulb-puzzle-vectors

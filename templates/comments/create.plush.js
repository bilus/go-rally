$('#comment-list').append('<%= partial("comments/item.plush.html") %>');
$('#comment-list-head').html('<%= partial("comments/head.html") %>')
Main.clearMarkdownEditor()

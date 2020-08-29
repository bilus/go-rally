$('#comment-list').append('<%= partial("comments/item.plush.html", { "comment": comment }) %>');
$('#comment-Body').val('');

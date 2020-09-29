$('#comment-list [data-comment-id="<%= comment.ID %>"]').remove()
$('#comment-list-head').html('<%= partial("comments/head.html") %>')

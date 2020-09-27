Main.hideTooltips()
$('.reactions[data-post-id="<%= post.ID %>"]').replaceWith('<%= partial("posts/reactions.html") %>')

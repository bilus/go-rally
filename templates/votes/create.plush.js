$('tr[data-post-id="<%= post.ID %>"]').replaceWith('<%= partial("posts/row.html", { "post": post }) %>');
$('.user-votes').replaceWith('<%= partial("user_votes.html") %>')

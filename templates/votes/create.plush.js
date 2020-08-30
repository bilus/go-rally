$('tr[data-post-id="<%= post.ID %>"]').replaceWith('<%= partial("posts/row.html", { "post": post }) %>');
$('.user-votes').html('<%= partial("user_votes.html") %>')

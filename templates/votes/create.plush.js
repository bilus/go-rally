$('tr[data-post-id="<%= post.ID %>"]').replaceWith('<%= partial("boards/post.html", { "post": post }) %>');
$('.user-votes').html('<%= partial("user_votes.html") %>')

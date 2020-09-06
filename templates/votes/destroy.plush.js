$('*[data-post-id="<%= post.ID %>"] .num-votes').replaceWith('<%= partial("boards/post_votes.html", { "post": post }) %>');
$('.user-votes').html('<%= partial("user_votes.html") %>')

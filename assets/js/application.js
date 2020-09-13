require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
import { EmojiButton } from '@joeattardi/emoji-button';

$(() => {
    $('[data-toggle="tooltip"]').tooltip()

    $('.clickable').click(function(e) {
        var tr = $(e.target).closest(".clickable");
        window.location = tr.data("href");
    })

    var editor = document.getElementById('post-Body');
    if (editor !== undefined) {
        var emojiPicker = new EmojiButton();
        var easyMDE = new EasyMDE({
            element: editor,
            forceSync: true,
            promptURLs: true,
            spellChecker: false,
            uploadImage: true,
            imageUploadEndpoint: "<%= postImagesPath({post_id: post.ID}) %>",
            imageCSRFToken: "<%= authenticity_token %>",
            toolbar: [
                'undo', 'redo',
                '|',
                'bold', 'italic', 'strikethrough', 'heading',
                '|',
                'code', 'quote', 'unordered-list', 'ordered-list',
                '|',
                'link', 'image',
                '|',
                'table', 'horizontal-rule',
                '|',
                'preview', 'side-by-side', 'fullscreen',
                '|',
                {
                    name: "emoji",
                    action: _editor => {
                        emojiPicker.togglePicker(document.querySelector(".EasyMDEContainer"))
                    },
                    className: "icon-star",
                    title: "Insert emoji",
                },
                "|",
                'guide',
            ]
        });

        emojiPicker.on('emoji', selection => {
            easyMDE.codemirror.replaceSelection(selection.emoji);
        });

        window.easyMDE = easyMDE;
    }
});

require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
var sentinel = require("sentinel-js");

import { EmojiButton } from '@joeattardi/emoji-button';

$(() => {
    $('[data-toggle="tooltip"]').tooltip()

    $('.clickable').click(e => {
        const tr = $(e.target).closest(".clickable");
        window.location = tr.data("href");
    })


    function installMarkdownEditor(editor) {
        const emojiPicker = new EmojiButton();
        const easyMDE = new EasyMDE({
            element: editor,
            forceSync: true,
            promptURLs: true,
            spellChecker: false,
            uploadImage: true,
            imageUploadEndpoint: $(editor).data("image-upload-endpoint"),
            // imageCSRFToken: "<%= authenticity_token %>",
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
    }

    sentinel.on('.markdown-editor', function(el) {
        installMarkdownEditor(el);
    });
    var editors = $(".markdown-editor");
    for (var i = 0; i < editors.length; i++) {
        const editor = editors[i];
        installMarkdownEditor(editor);
    }
});

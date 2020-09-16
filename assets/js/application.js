require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
var sentinel = require("sentinel-js");

var emoji = require('@joeattardi/emoji-button');
var markdownEditors = [];

$(() => {
    $('[data-toggle="tooltip"]').tooltip()

    $('.clickable').click(e => {
        const tr = $(e.target).closest(".clickable");
        window.location = tr.data("href");
    })


    function installMarkdownEditor(editor) {
        const editorHeight = $(editor).data("easymde-height");
        const emojiPicker = new emoji.EmojiButton();
        const easyMDE = new EasyMDE({
            element: editor,
            forceSync: true,
            promptURLs: true,
            spellChecker: false,
            uploadImage: true,
            minHeight: editorHeight,
            maxHeight: editorHeight,
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
                        emojiPicker.togglePicker(document.querySelector(".emoji"))
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
        return easyMDE;
    }

    sentinel.on('.markdown-editor', function(el) {
        markdownEditors.push(installMarkdownEditor(el));
    });
    var els = $(".markdown-editor");
    for (var i = 0; i < els.length; i++) {
        const el = els[i];
        markdownEditors.push(installMarkdownEditor(el));
    }
});

module.exports = {
    clearMarkdownEditor: function() {
        markdownEditors[0].codemirror.setValue('');
    }
};

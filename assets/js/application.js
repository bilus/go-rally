require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
var ClipboardJS = require("clipboard");
var sentinel = require("sentinel-js");

var emoji = require('@joeattardi/emoji-button');

var markdownEditors = [];

$(() => {
    hljs.configure({ languageDetectRe: '^highlight-.+$' });

    $('[data-toggle="tooltip"]').tooltip({ html: true, animation: false });
    sentinel.on('[data-toggle="tooltip"]', function(el) {
        $(el).tooltip({ html: true, animation: false });
    })

    $('[data-toggle="popover"]').popover({ html: true, animation: false })
    sentinel.on('[data-toggle="popover"]', function(el) {
        $(el).popover({ html: true, animation: false });
    })

    sentinel.on('.copy-board', function(el) {
        new ClipboardJS(el,
            {
                text: function(trigger) {
                    return window.location;
                }
            });
    })

    $('.clickable').click(e => {
        const tr = $(e.target).closest(".clickable");
        window.location = tr.data("href");
    })


    function installMarkdownEditor(editor) {
        const editorHeight = $(editor).data("easymde-height");
        const emojiPicker = new emoji.EmojiButton({ showAnimation: false });
        const easyMDE = new EasyMDE({
            element: editor,
            forceSync: true,
            promptURLs: true,
            spellChecker: false,
            uploadImage: true,
            minHeight: editorHeight,
            maxHeight: editorHeight,
            imageUploadEndpoint: $(editor).data("image-upload-endpoint"),
            renderingConfig: {
                codeSyntaxHighlighting: true,
            },
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
                    className: "icon-smile-o",
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

    sentinel.on('div.highlight pre', function(el) {
        hljs.highlightBlock(el);
    });
    document.querySelectorAll('div.highlight pre').forEach((el) => {
        hljs.highlightBlock(el);
    });


});

module.exports = {
    clearMarkdownEditor: function() {
        markdownEditors[0].codemirror.setValue('');
    },
    renderMarkdown: function(md) {
        var em = new EasyMDE({
            renderingConfig: {
                codeSyntaxHighlighting: true,
            },
        });
        return em.markdown(md);
    },
    remoteRequest: function(method, url, beforeSend) {
        $.ajax({
            method: method,
            url: url,
            dataType: 'script',
            contentType: 'text/javascript',
            beforeSend: beforeSend,
        });
    },
    pickEmoji: function(targetEl, callback) {
        const emojiPicker = new emoji.EmojiButton({
            showAnimation: false,
        });
        emojiPicker.on('emoji', callback)
        emojiPicker.togglePicker(targetEl)
    },
    onClick(sel, callback) {
        $(sel).click(callback);
        sentinel.on(sel, function(el) {
            $(el).click(callback)
        })
    },
    hideTooltips() {
        $('[data-toggle="tooltip"]').tooltip('hide');
    }
};

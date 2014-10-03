// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themefiles

const PresentationJs = `
$(function() {

	var presentationSelector = 'article.presentation > .content';

	var transformPresentationStructure = function() {
		var presentationContent = $(presentationSelector).html();
		var slides = presentationContent.split("<hr>")
		var newHtml = '<section class="slide">' + slides.join('</section><section class="slide">') + '</section>';
		$(presentationSelector).html(newHtml);
	};

	var renderPresentation = function() {


		if ($(presentationSelector).length == 0) {
			// this document is not a presentation
			return;
		}

		/**
		 * Toggle the page header elements
		 */
		var togglePresentationMode = function() {
			$("body>nav.toplevel").toggle();
			$("body>nav.breadcrumb").toggle();
			$("article.presentation>header").toggle();
			$("article.presentation>.description").toggle();
			$("body>footer").toggle();
		};

		// render the presentation
		$.deck('.slide', {
			selectors: {
				container: presentationSelector
			},

			keys: {
				goto: 71 // 'g'
			}
		});

		// handle keyboard shortcuts
		$(document).keydown(function(e) {

			/* <ctrl> + <shift> */
			if (e.ctrlKey && (e.which === 16) ) {
				console.log( "You pressed Ctrl + Shift" );
				togglePresentationMode();
			}

		});

	};

	// transform the content
	transformPresentationStructure();

	// render the presentaton
	renderPresentation();

    // register a on change listener
    if (typeof(autoupdate) === 'object' && typeof(autoupdate.onchange) === 'function') {
        autoupdate.onchange(
            "Render Presentation",
            function() {
                renderPresentation();
            }
        );
    }

});
`

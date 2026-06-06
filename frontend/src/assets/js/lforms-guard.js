// Externalized from index.html so a strict CSP can use script-src 'self' (#105/#103).
// Must stay a parser-blocking <script src> at the same position, immediately before the
// "closing marker" HTML comment in index.html: if the browser lacks Custom Elements,
// document.write('<!--') opens a comment that the following marker closes, disabling the
// lforms web-component loader. (lhncbc/lforms requirement.)
if (!window.customElements) {
  document.write('<!--');
}

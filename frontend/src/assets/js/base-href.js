// Externalized from index.html so a strict CSP can use script-src 'self' with no
// 'unsafe-inline' (#105/#103). Must stay a parser-blocking <script src> in <head>
// (no defer/async) so document.write runs during parse, before Angular reads <base>.
(function () {
  var baseHref = "/";

  // if the pathname includes /web, everything before `/web` (and including web) should be set as the base path.
  if (window.location.pathname.includes('/web')) {
    baseHref = "/web/";
    // probably running locally, and *may* include a subpath
    var subPath = window.location.pathname.split('/web').slice(0, 1)[0];
    if (subPath != "/") {
      // subpath, so we need to update the absolutePath with the subpath before adding the relative path to the end
      baseHref = subPath + '/web/';
    }
  }

  document.write('<base href="' + baseHref + '"/>');
})();

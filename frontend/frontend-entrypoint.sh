#!/usr/bin/env sh
# generate a tiny env.js the app can read at runtime
cat <<EOF > /usr/share/nginx/html/env.js
window.__env = {
  ALLOW_SIGNUP: "${ALLOW_SIGNUP}"
};
EOF

# launch nginx
exec nginx -g 'daemon off;'

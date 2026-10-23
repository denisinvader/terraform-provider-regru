Target nameserver domain. The provider normalizes it as follows:

    * If the value ends with `.`, it is passed unchanged:
        `ns1.example.com.` → `ns1.example.com.`.
    * If the value contains a dot but does not end with one, a trailing `.` is added:
        `ns1.example.com` → `ns1.example.com.`.
    * If the value contains no dot, `.<zone>.` is appended:
        `ns1` → `ns1.example.com.` when `zone` is `example.com`.

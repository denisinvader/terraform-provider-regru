Target mail server domain. The provider normalizes it as follows:

    * If the value ends with `.`, it is passed unchanged:
        `mx.example.com.` → `mx.example.com.`.
    * If the value contains a dot but does not end with one, a trailing `.` is added:
        `mx.example.com` → `mx.example.com.`.
    * If the value contains no dot, `.<zone>.` is appended:
        `mx` → `mx.example.com.` when `zone` is `example.com`.

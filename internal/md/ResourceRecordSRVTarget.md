Value passed to reg.ru as `target`. The provider normalizes it as follows:

    * If the value ends with `.`, it is passed unchanged:
        `sip.example.com.` → `sip.example.com.`.
    * If the value contains a dot but does not end with one, a trailing `.` is added:
        `sip.example.com` → `sip.example.com.`.
    * If the value contains no dot, `.<zone>.` is appended:
        `sip` → `sip.example.com.` when `zone` is `example.com`.

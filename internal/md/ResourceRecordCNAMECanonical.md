Target domain. The provider normalizes it as follows:

    * If the value ends with `.`, it is passed unchanged:
        `example.com.` → `example.com.`.
    * If the value contains a dot but does not end with one, a trailing `.` is added:
        `foo.example.com` → `foo.example.com.`.
    * If the value contains no dot, `.<zone>.` is appended:
        `foo-bar` → `foo-bar.example.com.` when `zone` is `example.com`.

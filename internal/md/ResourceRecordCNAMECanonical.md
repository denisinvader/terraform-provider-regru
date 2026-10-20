Target domain name.

Following normalization rules apply:

* If value ends with `.` it used as is:
    `example.com. -> example.com.`.
* Else if value contains`.` in the middle, `.` concatenates to the end:
    `foo.example.com -> foo.example.com.`.
* Else zone concatenates to the end:
    `foo-bar -> foo-bar.example.com.`.

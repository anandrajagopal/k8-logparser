# k8-logparser

This is a small utility that runs as a side car to parse logs from a main container. Tools like fluentd parse each line separately and this causes issues when an application runs inside a K8 pod and writes logs to stdout. For example, if java verbose gc and application logs (and potentially other kinds of logs) go to stdout, it is difficult to know how many lines of logs need to be grouped together before logs can be shipped to ELK. This tool/toy program/proof of concept attempts to solve it by having the application write the logs to a named pipe. Each type of log would be written to a named pipe and this side car container (it's configurable) reads from those pipes concurrently using the specified regex. The parser starts by reading from the named pipe and starts accumulating lines until the specified regex matches. Once the regex matches, the parser adds '{{' to the beginning and '}}' to the end of the accumulated lines and pushes it to stdout. 
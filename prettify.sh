#!/usr/bin/env python3

import sys
from html.parser import HTMLParser

class PrettyPrinter(HTMLParser):
    def __init__(self):
        super().__init__()
        self.indent = 0
        self.output = []

    def handle_starttag(self, tag, attrs):
        self.output.append('  ' * self.indent + self.get_starttag_text())
        if tag not in ['br', 'img', 'input', 'hr', 'meta', 'link']:
            self.indent += 1
        self.output.append('\n')

    def handle_endtag(self, tag):
        if tag not in ['br', 'img', 'input', 'hr', 'meta', 'link']:
            self.indent -= 1
            self.output.append('  ' * self.indent + f'</{tag}>\n')

    def handle_data(self, data):
        if data.strip():
            self.output.append('  ' * self.indent + data.strip() + '\n')

    def handle_comment(self, data):
        self.output.append('  ' * self.indent + f'<!--{data}-->\n')

def prettify_html(input_html):
    parser = PrettyPrinter()
    parser.feed(input_html)
    return ''.join(parser.output)

if __name__ == "__main__":
    input_html = sys.stdin.read()
    pretty_html = prettify_html(input_html)
    print(pretty_html)

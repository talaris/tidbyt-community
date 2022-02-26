"""
Applet: Kickstarter
Summary: Track Campaign
Description: Trakcs the progress of a Kickstarter project.
Author: talaris
"""

load("render.star", "render")
load("schema.star", "schema")

def main(config):
    return render.Root(
        child = render.Box(height = 1, width = 1),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [],
    )
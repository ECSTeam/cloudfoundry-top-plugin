# top-plugin Screenshots


### Application view:
Shows all applications deployed to foundation. Including app CPU and memory usage.
This is useful for finding the "top" CPU consumer on the foundation.  This screenshot 
also shows an alert:
![Screenshot](screenshot_appView2.png?raw=true)


### Application view - multi-stack:
Top supports multiple stacks.  This screenshot shows a foundation that has both cflinuxfs2 and windows2012R2 cells.
Top lets you filter output which enable you to focus in on the applications that are of interest.
![Screenshot](screenshot_appViewMultiStack.png?raw=true)

### Application view - filters:
Top lets you filter output which enable you to focus in on the applications that are of interest.  
All columns can be used for filtering.  Alphanumeric columns support regular expressions.  Numeric columns allow
simple expressions (e.g., >15 ).  Filters on multiple columns are treated as an "and" condition.
![Screenshot](screenshot_appViewFilter.png?raw=true)


### Application view - sorting:
Top allows sorting output on any column. Sorting on multiple columns is supporting (up to 5 levels).
![Screenshot](screenshot_appViewSort.png?raw=true)


### Application "detail" view:
Shows all instances (containers) of selected application.
![Screenshot](screenshot_appDetailView.png?raw=true)


### Application "detail" - info view:
Shows additional information about selected application.
![Screenshot](screenshot_appDetailViewAppInfo.png?raw=true)


### Diego Cell view:
Shows all Diego cells running on foundation. Includes cell CPU and memory usage. 
This view is useful for locating any "hot" cell -- a cell that has a higher then
expected CPU utilization.
![Screenshot](screenshot_cellView.png?raw=true)


### Diego Cell "Detail" view:
Shows all containers running on selected cell. Includes container CPU and memory usage.  
If a cell has a high CPU utilization, this detail view can help identify which application instance is the culprit. 
![Screenshot](screenshot_cellDetailView.png?raw=true)


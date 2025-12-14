# TODO

- BUG: Position of CustomSelect hover goes outside the screen for task at the bottom of the List.
- BUG: Very High Tab scrollbar does not work.
  (issue is that we use "CriticalTasks" for counting the elements for all Tabs!)
- Add some visual cue in the Task list when a Task has a Note.
- BUG: Using td.CriticalTasks in task getGenericTaskTable!!! need to use the currently selected Tab matching Array instead!!
- Add attachment to notes or as a separate functionality
- Change customSelect "hover" color to something a little nicer
- Add refresh to Priority change as well (just like status)
- Create a new branch for Position Drag and Drop
- IMPROVEMENT: Need to improve the up/down position... it's hard to use and the Position Number does not adjust properly when task are set to "Done".
- Auto-refresh the updated_at field
- Split task in two rows and add timer / Total time
- Add user info in DB for the future maybe?
- BUG: Sometime auto-switching the Status to not started in UI when clicking in the task row
- BUG: Position UP/DOWN seems off... probably related to switching some task to "Done"
- Need to be able to set the default Value / Color on CustomSelect
- Delete task?
- Limit the length of the Title field
- Add comments / update


# DONE

- (/) Add: Notes listing should be scrollable.
- (/) Set focus to field in OnSelected Form (edit title) automatically
- (/) Disable Up button on first row, disable down button in last row
- (/) Popup fonts looks blur and does not follow the color of the other Text
- (/) Add tab switchting with cmd+number
- (/) Remove empty option in customSelect for both Status and Priority
- (/) Make the CustomSelect text bold in popup menu
- (/) Hide what is not working for the moment (all the files functionalities)
- (/) Show Done Tasks somewhere
- (/) BUG: Prevent addition of "empty" Task (clicking multiple time the Add new Task Button) 
- (/) BUG: Cannot select text in notes
    Now the whole text can be copied with the copy button
- (/) Replace the status switch origTab with refreshStatusTab()
- (/) BUG: Adding a new Task should refresh the current Tab
- (/) BUG: In TODO Tab, status changes do not trigger view update
- (/) BUG: Very High, High- Tab shows ghost value on adding / removing (set status to done) task
- (/) Add all other Priorities Tabs as well
- (/) BUG: Up/Down position jammed after first 
- (/) BUG: Switching task to "Done" does not refresh inactive Tab
- (/) BUG: Switching Status does not affect other tab and when switching back to orignal Tab... status are reverted back (only in UI, DB is fine)moved... need to update other Tab as well
- (/) BUG: Table refresh not working as before. When adding new task, it does not appears in the table automatically.
- (/) Get urgently a 'critical' Tab!
- (/) Improve update_at formating
- (/) Add Search
- (/) Bundle placeholder image
- (/) Edit title on click
- (/) When switching task to done... need to update all position
- (/) Replace position field for a relation table Task_Status_Position. Giving the possibility to order per Status. Important is to get position 'unique' for "non" completed tasks for the main view
- (/) Add DB filename filter (*.medo)
- (/) Replace dropdown for two button for UP / Down position
- (/) Change Custom-Select font to white
- (/) Adjust the size of the CustomSelect automatically depending on the content
- (/) Need bigger placeholder image to increase the size of the Window
- (/) Get the scrollbar to work in task list (convert to Table??)
- (/) Update UI table after new task is being added to DB
- (/) Adding a task now crash the App
- (/) 'Sample task' are being add automatically even if the table is not empty. That should not happen.
- (/) Refactor task grid for auto-refresh when changing position
- (/) Add all Priority dropdown values
- (/)Find a way to "enlarge" the titiel field in new task form


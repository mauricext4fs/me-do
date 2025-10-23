# TODO

- BUG: Cannot select text in notes
- Create a new branch for Position Drag and Drop
- IMPROVEMENT: Need to improve the up/down position... it's hard to use and the Position Number does not adjust properly when task are set to "Done".
- Auto-refresh the updated_at field
- Split task in two rows and add timer / Total time
- Add user info in DB for the future maybe?
- BUG: Sometime auto-switching the Status to not started in UI when clicking in the task row
- BUG: Position UP/DOWN seems off... probably related to switching some task to "Done"
- Disable Up button on first row, disable down button in last row
- Set focus to field in OnSelected Form automatically
- Need to be able to set the default Value / Color on CustomSelect
- Make the CustomSelect text bold in popup menu
- Status and Priority background color should also show in Table
- Delete task?
- Limit the length of the Title field
- Split view for "Done" vs not done tasks
- Add comments / update


# DONE

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


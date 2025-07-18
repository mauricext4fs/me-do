# TODO

- Need bigger placeholder image to increase the size of the Window
- Need to be able to set the default Value / Color on CustomSelect
- Adjust the dropdown color
- Adjust the size of the CustomSelect automatically depending on the content
- Make the CustomSelect text bold in popup menu
- Status and Priority background color should also show in Table
- Delete task?
- Limit the length of the Title field
- Add event on changing status to "done" to update position (update task set posistion = position -1)
- Add event on changing "done" status to reorder position (update task set position = position + 1), update task set position = 1 where id = xxx
- Split view for "Done" vs not done tasks
- Replace position field for a relation table Task_Status_Position. Giving the possibility to order per Status. Important is to get position 'unique' for "non" completed tasks for the main view
- Replace dropdown for two button for UP / Down position
- Add comments / update


# DONE

- (/) Get the scrollbar to work in task list (convert to Table??)
- (/) Update UI table after new task is being added to DB
- (/) Adding a task now crash the App
- (/) 'Sample task' are being add automatically even if the table is not empty. That should not happen.
- (/) Refactor task grid for auto-refresh when changing position
- (/) Add all Priority dropdown values
- (/)Find a way to "enlarge" the titiel field in new task form


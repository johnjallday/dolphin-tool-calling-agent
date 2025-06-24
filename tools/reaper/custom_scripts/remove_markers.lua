function main()
	-- Begin an undo block so the action can be reversed with a single undo.
	reaper.Undo_BeginBlock()

	local markers_removed_count = 0

	-- reaper.CountProjectMarkers() returns the total number of both markers and regions.
	-- We need to loop backwards when deleting, otherwise the list gets re-indexed
	-- and the script would skip over some markers.
	local num_items = reaper.CountProjectMarkers(0)

	for i = num_items - 1, 0, -1 do
		-- EnumProjectMarkers gets information about the marker/region at a given index.
		-- The key is the 'isrgn' return value, which is true for regions and false for markers.
		--
		reaper.DeleteProjectMarker(0, i, false)

		reaper.Undo_EndBlock("Remove All Markers (Excluding Regions)", -1)
	end
end

-- Defer the execution to ensure Reaper's UI is ready.
reaper.defer(function()
	main()
	reaper.UpdateArrange() -- Refresh the ruler to show the markers have been removed.
end)

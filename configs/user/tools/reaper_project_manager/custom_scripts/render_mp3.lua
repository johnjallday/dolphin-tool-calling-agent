function getMarkerPositionByName(name_to_find)
	local num_items = reaper.CountProjectMarkers(0)
	for i = 0, num_items - 1 do
		local retval, _, pos, _, name = reaper.EnumProjectMarkers3(0, i)
		if retval > 0 and name == name_to_find then
			return pos -- Return the position if found
		end
	end
	return nil -- Return nil if not found
end

function render_project_between_markers()
	-- Find the start and end marker positions
	local start_pos = getMarkerPositionByName("start")
	local end_pos = getMarkerPositionByName("end")

	-- Check if markers were found
	if not start_pos or not end_pos then
		reaper.ShowConsoleMsg("Error: 'start' or 'end' marker not found. Cannot set render range.\n")
		return
	end

	if start_pos >= end_pos then
		reaper.ShowConsoleMsg("Error: 'start' marker is at or after 'end' marker. Cannot set render range.\n")
		return
	end

	--reaper.ShowConsoleMsg("Found 'start' at " .. start_pos .. " and 'end' at " .. end_pos .. ".\n")

	-- Set the time selection to the marker range
	reaper.GetSet_LoopTimeRange2(0, true, false, start_pos, end_pos, false)

	-- Set render settings
	reaper.GetSetProjectInfo(0, "RENDER_TIMESOURCE", 2, true) -- Set render source to "Time Selection"
	reaper.GetSetProjectInfo(0, "RENDER_SRATE", 48000, true) -- Set sample rate to 48000 Hz
	reaper.GetSetProjectInfo_String(0, "RENDER_FORMAT", "l3pm", true) -- Set format to MP3 (LAME)
	reaper.GetSetProjectInfo_String(0, "RENDER_PATTERN", "$project", true) -- Set output file pattern to project name

	-- Open the Render to File dialog
	-- Action ID 42230 is for "File: Render project..."
	reaper.Main_OnCommand(42230, 0)
end

----------------------------------------------------------------------

render_project_between_markers()

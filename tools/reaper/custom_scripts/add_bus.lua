-- description
-- add BUS tracks on the top

function findTrackByName(name)
	for i = 0, reaper.CountTracks(0) - 1 do
		local track = reaper.GetTrack(0, i)
		-- Get the name of the current track
		local _, trackName = reaper.GetSetMediaTrackInfo_String(track, "P_NAME", "", false)
		if trackName == name then
			return track -- Return the track object if found
		end
	end
	return nil -- Return nil if no track with that name is found
end

function main()
	-- Define the names for the parent folder and its children
	local parent_track_name = "BUS MASTER"
	local child_track_names = { "BUS INST", "BUS VOX", "BUS SEND" }

	-- Check if the "BUS MASTER" track already exists before doing anything
	if findTrackByName(parent_track_name) then
		reaper.ShowConsoleMsg("'" .. parent_track_name .. "' track already exists. Script aborted.\n")
		return -- Exit the script if the track is found
	end

	-- If the track doesn't exist, proceed with creation
	reaper.Undo_BeginBlock()

	-- 1. Insert the main "BUS MASTER" track at the top (index 0)
	reaper.InsertTrackAtIndex(0, true)
	local parent_track = reaper.GetTrack(0, 0)
	reaper.GetSetMediaTrackInfo_String(parent_track, "P_NAME", parent_track_name, true)
	reaper.SetTrackSelected(parent_track, true) -- Select the parent track

	-- 2. Insert the child tracks immediately after the parent track
	for i, name in ipairs(child_track_names) do
		reaper.InsertTrackAtIndex(i, true) -- Insert at index 1, then 2, then 3
		local child_track = reaper.GetTrack(0, i)
		reaper.GetSetMediaTrackInfo_String(child_track, "P_NAME", name, true)

		-- If the track is "BUS INST", set its volume to -6dB
		if name == "BUS INST" then
			reaper.SetMediaTrackInfo_Value(child_track, "D_VOL", 0.50)
		end

		reaper.SetTrackSelected(child_track, true) -- Select each child track as it's created
	end

	-- 3. Apply the folder settings to the correct tracks
	reaper.SetMediaTrackInfo_Value(parent_track, "I_FOLDERDEPTH", 1)
	local last_child_track_index = #child_track_names
	local last_child_track = reaper.GetTrack(0, last_child_track_index)
	reaper.SetMediaTrackInfo_Value(last_child_track, "I_FOLDERDEPTH", -1)

	-- 4. Minimize the height of all newly created (and now selected) tracks
	reaper.Main_OnCommand(40122, 0) -- Action: "Track: Set selected track(s) height to minimum"

	-- End the undo block
	reaper.Undo_EndBlock("Create Bus Master Folder and Set Volume", -1)
end

-- Defer the script's execution for stability
reaper.defer(function()
	main()
	reaper.UpdateArrange() -- Update the track view to show the changes
end)

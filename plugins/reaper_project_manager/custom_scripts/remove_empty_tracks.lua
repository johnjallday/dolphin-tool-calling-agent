--remove empty tracks
--exceptions BUS_routed or parent folder
-- ReaScript to remove empty tracks with exceptions for BUS routed or parent folder tracks

-- --- Constants ---
local DEBUG = false -- Set to false to disable debug messages in the Reaper console

-- --- Helper Functions ---

function IsTrackRoutedToBus(track)
	-- Checks if a track has any sends to another track (acting as a bus).
	local num_sends = reaper.GetTrackNumSends(track, 0) -- category 0 for track-to-track sends
	return num_sends > 0
end

function IsTrackInFolder(track)
	-- Checks if a track is a child in a folder.
	local parent_track = reaper.GetParentTrack(track)
	if parent_track then
		-- Check if the parent is actually a folder start
		-- A track is a folder start if its depth is positive and the next track's depth is greater
		local parent_depth = reaper.GetTrackDepth(parent_track)
		local parent_idx = reaper.GetMediaTrackInfo_Value(parent_track, "IP_TRACKNUMBER") - 1 -- 0-indexed

		if parent_idx < reaper.CountTracks(0) - 1 then -- Check if it's not the last track
			local next_track = reaper.GetTrack(0, parent_idx + 1)
			if next_track then -- Ensure next_track exists
				local next_track_depth = reaper.GetTrackDepth(next_track)
				if parent_depth >= 0 and next_track_depth > parent_depth then
					return true
				end
			end
		end
	end
	return false
end

function IsTrackEmpty(track)
	-- Checks if a track has no media items and no FX.
	local num_items = reaper.CountTrackMediaItems(track)
	local num_fx = reaper.TrackFX_GetCount(track)
	return num_items == 0 and num_fx == 0
end

function GetTrackName(track)
	-- Gets the name of a track.
	local retval, _, name_buffer = reaper.GetTrackName(track, "", 256) -- Provide a buffer
	if retval then
		return name_buffer
	end
	return "[Unknown Name]"
end

-- --- Main Script Logic ---

function RemoveEmptyTracksWithExceptions()
	reaper.Undo_BeginBlock() -- Start of undo block for a single undo point

	local tracks_to_remove = {}
	local num_tracks = reaper.CountTracks(0) -- 0 for current project

	if DEBUG then
		reaper.ShowConsoleMsg("Starting script: Remove Empty Tracks (Lua)...\n")
		reaper.ShowConsoleMsg("Total tracks found: " .. num_tracks .. "\n")
	end

	-- Iterate backwards because deleting tracks shifts indices
	for i = num_tracks - 1, 0, -1 do
		local track = reaper.GetTrack(0, i)
		if track then -- Ensure track exists
			local track_name = GetTrackName(track)

			if DEBUG then
				reaper.ShowConsoleMsg("\nProcessing Track #" .. (i + 1) .. ": " .. track_name .. "\n")
			end

			if IsTrackEmpty(track) then
				if DEBUG then
					reaper.ShowConsoleMsg("  Track is empty (no items, no FX).\n")
				end

				local is_routed_to_bus = IsTrackRoutedToBus(track)
				local is_in_folder = IsTrackInFolder(track)

				if is_routed_to_bus then
					if DEBUG then
						reaper.ShowConsoleMsg(
							"  EXCEPTION: Track '" .. track_name .. "' is routed to a BUS. Will not remove.\n"
						)
					end
				-- continue equivalent in Lua is to just let the loop proceed
				elseif is_in_folder then
					if DEBUG then
						reaper.ShowConsoleMsg(
							"  EXCEPTION: Track '" .. track_name .. "' is part of a parent folder. Will not remove.\n"
						)
					end
				-- continue equivalent
				else
					-- If no exceptions apply, mark for removal
					if DEBUG then
						reaper.ShowConsoleMsg(
							"  Track '" .. track_name .. "' is empty and meets no exceptions. Marking for removal.\n"
						)
					end
					table.insert(tracks_to_remove, track)
				end
			else
				if DEBUG then
					reaper.ShowConsoleMsg("  Track is not empty.\n")
				end
			end
		end
	end

	-- Actually delete the marked tracks
	if #tracks_to_remove > 0 then
		if DEBUG then
			reaper.ShowConsoleMsg("\nTracks to be removed: " .. #tracks_to_remove .. "\n")
		end
		for _, track_to_delete in ipairs(tracks_to_remove) do
			local track_name_deleted = GetTrackName(track_to_delete) -- Get name again
			reaper.DeleteTrack(track_to_delete)
			if DEBUG then
				reaper.ShowConsoleMsg("  Removed track: " .. track_name_deleted .. "\n")
			end
		end
		--reaper.ShowConsoleMsg("Successfully removed " .. #tracks_to_remove .. " empty track(s).\n")
	else
		if DEBUG then
			reaper.ShowConsoleMsg("No empty tracks found matching removal criteria.\n")
		end
	end

	reaper.Undo_EndBlock("Remove Empty Tracks (Lua, with exceptions)", -1) -- End of undo block
	if DEBUG then
		reaper.ShowConsoleMsg("Script finished.\n")
	end
	-- Optional: reaper.UpdateArrange() if you want to visually refresh the track view immediately,
	-- though Reaper usually handles this.
end

RemoveEmptyTracksWithExceptions()

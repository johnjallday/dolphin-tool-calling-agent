-- ReaScript: add_kick_4_on_the_floor
-- Author: Gemini
-- Version: 1.2 (Lua)
-- Description: Creates a new track with ReaSamplOmatic5000 and a four-on-the-floor kick pattern.
-- Fixes:
-- v1.1: Used GetSetMediaTrackInfo_String to correctly set the track name.
-- v1.2: Used GetProjectTimeSignature for more robustly fetching beats per measure.

function main()
	reaper.Undo_BeginBlock()

	-- Create a new track
	local track_index = reaper.CountTracks()
	reaper.InsertTrackAtIndex(track_index, true)
	local new_track = reaper.GetTrack(0, track_index)

	-- Name the track "Kick" using the correct function for strings
	reaper.GetSetMediaTrackInfo_String(new_track, "P_NAME", "Kick", true)

	-- Add ReaSamplOmatic5000 to the new track
	local vst_name = "ReaSamplOmatic5000 (Cockos)"
	local fx_index = reaper.TrackFX_AddByName(new_track, vst_name, false, -1)

	if fx_index ~= -1 then
		-- Get the number of beats per measure from the project time signature
		local _, beats_per_measure = reaper.GetProjectTimeSignature(0)

		-- Check if we got a valid number, otherwise default to 4
		if not beats_per_measure or beats_per_measure == 0 then
			beats_per_measure = 4
		end

		-- Create a new MIDI item with the length of one measure (in quarter notes)
		local item = reaper.CreateNewMIDIItemInProj(new_track, 0, beats_per_measure, false)

		if item then
			local take = reaper.GetActiveTake(item)

			-- Define MIDI note properties for the kick drum
			local kick_note = 36 -- C1, a common MIDI note for kick drums
			local velocity = 100
			local ppq = 960 -- Standard Pulses Per Quarter note

			-- Insert four kick notes on each beat of the measure
			for i = 0, 3 do
				local position = i * ppq
				local length = ppq / 4 -- 16th note duration
				reaper.MIDI_InsertNote(take, false, false, position, position + length, 0, kick_note, velocity, false)
			end

			reaper.UpdateArrange()
		end

		reaper.Undo_EndBlock("Add Kick 4 on the Floor Track", -1)
	else
		reaper.ShowConsoleMsg("Failed to add ReaSamplOmatic5000. Make sure the plugin is available.\n")
		reaper.Undo_EndBlock("Attempted to Add Kick Track", -1)
	end
end

-- Check for a valid Reaper environment before running
if reaper then
	main()
else
	reaper.ShowConsoleMsg("This script must be run from within REAPER.\n")
end

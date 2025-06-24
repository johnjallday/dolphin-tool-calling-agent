function add_guitar()
	reaper.InsertTrackAtIndex(reaper.CountTracks(0), true)
	local track = reaper.GetTrack(0, reaper.CountTracks(0) - 1)
	reaper.GetSetMediaTrackInfo_String(track, "P_NAME", "M Guitar", true)
	plugin_name = "AU: Guitar Rig 7 FX (Native Instruments)"
	local ret, insert_fx_idx = reaper.TrackFX_AddByName(track, plugin_name, false, 1)
	reaper.SetMediaTrackInfo_Value(track, "I_RECARM", 1)
end

add_guitar()

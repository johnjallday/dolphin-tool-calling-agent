-- rearrange_tracks.lua (Refactored with Workaround for InsertTrackAtIndex)

-- User-configurable constants
local DEBUG = true -- Set to false to disable most console messages

-- Target track names
local BUS_MASTER_NAME = "BUS MASTER"
local BUS_INST_NAME = "BUS INST"
local BUS_VOX_NAME = "BUS VOX"

local GROUP_DRUMS_NAME = "GROUP DRUMS"
local GROUP_BASS_NAME = "GROUP BASS"
local GROUP_MELODY_NAME = "GROUP MELODY"
local GROUP_FX_NAME = "GROUP FX"

-- Prefixes for grouping (ensure the space is included if that's the convention)
local PREFIX_DRUMS = "DR "
local PREFIX_BASS = "BA "
local PREFIX_MELODY = "M "
local PREFIX_FX = "FX "

-------------------------------------------------------------------------------
-- HELPER FUNCTIONS
-------------------------------------------------------------------------------

function GetTrackByName(name_to_find)
  if DEBUG then reaper.ShowConsoleMsg("Fn GetTrackByName: Searching for track: '" .. name_to_find .. "'\n") end
  for i = 0, reaper.CountTracks(0) - 1 do
    local track = reaper.GetTrack(0, i)
    if track then
	 local retval, _, current_name = reaper.GetTrackName(track, "", 256)
	 if retval and current_name == name_to_find then
	   if DEBUG then reaper.ShowConsoleMsg("  Fn GetTrackByName: Found track: '" .. name_to_find .. "'\n") end
	   return track
	 end
    end
  end
  if DEBUG then reaper.ShowConsoleMsg("  Fn GetTrackByName: Track NOT found: '" .. name_to_find .. "'\n") end
  return nil
end

---
-- Creates a new track with a specific name at a desired index.
-- Includes a workaround for InsertTrackAtIndex returning nil when a track is actually created.
-- @param name (string) The desired name for the new track.
-- @param desired_idx (number) The 0-based index where the track should be inserted.
-- @return MediaTrack* The created track object, or nil on failure.
---
function create_track(name, desired_idx)
  if DEBUG then reaper.ShowConsoleMsg("Fn create_track: Attempting to create NEW track '" .. name .. "' at index " .. desired_idx .. "\n") end

  if type(name) ~= "string" or name == "" then
    if DEBUG then reaper.ShowConsoleMsg("  ERROR (create_track): Invalid track name provided (nil, empty, or not a string).\n") end
    return nil
  end
  if type(desired_idx) ~= "number" then
    if DEBUG then reaper.ShowConsoleMsg("  ERROR (create_track): desired_idx is not a number. Got: " .. tostring(desired_idx) .. "\n") end
    return nil
  end
  if desired_idx < 0 or desired_idx > reaper.CountTracks(0) then
	if DEBUG then reaper.ShowConsoleMsg("  ERROR (create_track): desired_idx " .. desired_idx .. " is out of valid range [0, " .. reaper.CountTracks(0) .. "]. Cannot insert.\n") end
	return nil
  end

  local count_before_insert = reaper.CountTracks(0)
  local returned_idx_from_api = reaper.InsertTrackAtIndex(desired_idx, false) -- make_selected = false
  local count_after_insert = reaper.CountTracks(0)

  local track_object_to_use = nil
  local actual_idx_of_new_track = -1 -- Initialize to invalid state

  if type(returned_idx_from_api) == "number" and returned_idx_from_api >= 0 then
    if DEBUG then reaper.ShowConsoleMsg("  (create_track) InsertTrackAtIndex API returned valid index: " .. returned_idx_from_api .. "\n") end
    actual_idx_of_new_track = returned_idx_from_api
    track_object_to_use = reaper.GetTrack(0, actual_idx_of_new_track)
  elseif count_after_insert == count_before_insert + 1 then
    -- WORKAROUND: API returned nil/error, but track count increased.
    -- This implies the track was inserted at desired_idx.
    if DEBUG then
	   reaper.ShowConsoleMsg("  (create_track) InsertTrackAtIndex API returned: " .. tostring(returned_idx_from_api) .. " (unexpected or error value).\n")
	   reaper.ShowConsoleMsg("  (create_track) BUT track count increased from " .. count_before_insert .. " to " .. count_after_insert .. ".\n")
	   reaper.ShowConsoleMsg("  (create_track) WORKAROUND ACTIVATED: Assuming track was inserted at desired_idx: " .. desired_idx .. "\n")
    end
    actual_idx_of_new_track = desired_idx -- The index where we asked it to be put
    track_object_to_use = reaper.GetTrack(0, actual_idx_of_new_track)
    if not track_object_to_use and DEBUG then
	   reaper.ShowConsoleMsg("  (create_track) WORKAROUND FAILED to get track: reaper.GetTrack(0, " .. actual_idx_of_new_track .. ") returned nil.\n")
    end
  else
    -- API returned error AND track count did not increase as expected for a new track.
    if DEBUG then
	   reaper.ShowConsoleMsg("  ERROR (create_track): InsertTrackAtIndex API returned: " .. tostring(returned_idx_from_api) .. ".\n")
	   reaper.ShowConsoleMsg("  (create_track) Track count changed from " .. count_before_insert .. " to " .. count_after_insert .. " (not an increase of 1 or same on error).\n")
	   reaper.ShowConsoleMsg("  (create_track) Concluding track insertion failed for '" .. name .. "'.\n")
    end
    return nil -- Definite failure
  end

  if not track_object_to_use then
    if DEBUG then reaper.ShowConsoleMsg("  ERROR (create_track): Failed to get track object for '" .. name .. "' (inferred/actual index was " .. actual_idx_of_new_track .. ").\n") end
    return nil
  else
    local _, _, initial_name_after_creation = reaper.GetTrackName(track_object_to_use, "", 256)
    if DEBUG then
	   reaper.ShowConsoleMsg("  (create_track) Retrieved track object. Initial default name: '" .. (initial_name_after_creation or "N/A") .. "'\n")
    end

    if DEBUG then reaper.ShowConsoleMsg("  (create_track) Attempting to set name to: '" .. name .. "'\n") end
    local success_rename = reaper.GetSetMediaTrackInfo_String(track_object_to_use, "P_NAME", name, true)
    local _, _, name_after_attempt = reaper.GetTrackName(track_object_to_use, "", 256)

    if DEBUG then
	   if success_rename then
		  reaper.ShowConsoleMsg("    (create_track) API call GetSetMediaTrackInfo_String reported SUCCESS.\n")
		  reaper.ShowConsoleMsg("    (create_track) Track name is NOW: '" .. (name_after_attempt or "N/A") .. "'\n")
		  if name_after_attempt ~= name then
			 reaper.ShowConsoleMsg("    (create_track) WARNING: Name after reported success ('" .. (name_after_attempt or "N/A") .. "') does NOT match desired name ('" .. name .. "').\n")
		  end
	   else
		  reaper.ShowConsoleMsg("    (create_track) API call GetSetMediaTrackInfo_String reported FAILURE.\n")
		  reaper.ShowConsoleMsg("    (create_track) Track name is currently: '" .. (name_after_attempt or "N/A") .. "'\n")
	   end
    end

    if not success_rename or name_after_attempt ~= name then
	   if DEBUG then reaper.ShowConsoleMsg("  ERROR (create_track): Renaming was not successful for '".. name .."'. Final name: '".. (name_after_attempt or "N/A") .."'.\n") end
	   -- If we "created" a track (via workaround) but couldn't name it, it's problematic.
	   -- Consider this a failure for the script's purpose.
	   return nil
    end

    if DEBUG then reaper.ShowConsoleMsg("  (create_track) Successfully created (or inferred creation of) and named track: '" .. name .. "'.\n") end
    return track_object_to_use
  end
end

function CreateTrackIfNeeded(name_to_create)
  local existing_track = GetTrackByName(name_to_create)
  if existing_track then
    if DEBUG then reaper.ShowConsoleMsg("Fn CreateTrackIfNeeded: Track '" .. name_to_create .. "' already exists. Using existing.\n") end
    return existing_track
  else
    if DEBUG then reaper.ShowConsoleMsg("Fn CreateTrackIfNeeded: Track '" .. name_to_create .. "' not found. Calling create_track to make it at the end of the list.\n") end
    return create_track(name_to_create, reaper.CountTracks(0))
  end
end

function SendExists(src_track, dest_track)
    if not src_track or not dest_track then return false end
    local send_count = reaper.GetTrackNumSends(src_track, 0)
    for i = 0, send_count - 1 do
	   local p_dest_track = reaper.GetTrackSendInfo_Value(src_track, 0, i, "P_DESTTRACK")
	   if p_dest_track == dest_track then
		  return true
	   end
    end
    return false
end

-------------------------------------------------------------------------------
-- MAIN FUNCTION
-------------------------------------------------------------------------------
function RearrangeTracks()
  reaper.Undo_BeginBlock()
  reaper.PreventUIRefresh(1)

  if DEBUG then reaper.ShowConsoleMsg("Starting Track Rearrangement Script...\n") end

  if DEBUG then reaper.ShowConsoleMsg("\n--- Phase 1 & 2: Ensuring Core Tracks Exist ---\n") end
  local bus_master = CreateTrackIfNeeded(BUS_MASTER_NAME)
  local bus_inst = CreateTrackIfNeeded(BUS_INST_NAME)
  local bus_vox = CreateTrackIfNeeded(BUS_VOX_NAME)

  local group_drums = CreateTrackIfNeeded(GROUP_DRUMS_NAME)
  local group_bass = CreateTrackIfNeeded(GROUP_BASS_NAME)
  local group_melody = CreateTrackIfNeeded(GROUP_MELODY_NAME)
  local group_fx = CreateTrackIfNeeded(GROUP_FX_NAME)

  if not bus_master or not bus_inst or not bus_vox or not group_drums or not group_bass or not group_melody or not group_fx then
	 if DEBUG then reaper.ShowConsoleMsg("CRITICAL ERROR: One or more core tracks could not be created or found. Aborting script.\n") end
	 reaper.PreventUIRefresh(-1)
	 reaper.Undo_EndBlock("Rearrange Tracks (Failed - Core Track Creation)", -1)
	 reaper.UpdateArrange()
	 return
  end

  local core_tracks_map = {
    [bus_master] = BUS_MASTER_NAME, [bus_inst] = BUS_INST_NAME, [bus_vox] = BUS_VOX_NAME,
    [group_drums] = GROUP_DRUMS_NAME, [group_bass] = GROUP_BASS_NAME, [group_melody] = GROUP_MELODY_NAME, [group_fx] = GROUP_FX_NAME
  }

  if DEBUG then reaper.ShowConsoleMsg("\n--- Collecting Child Tracks ---\n") end
  local drum_children, bass_children, melody_children, fx_children, other_tracks = {}, {}, {}, {}, {}

  for i = 0, reaper.CountTracks(0) - 1 do
    local track = reaper.GetTrack(0, i)
    if track then
	   if not core_tracks_map[track] then
		local retval, _, name = reaper.GetTrackName(track, "", 256)
		if retval and name then
		  if string.find(name, PREFIX_DRUMS, 1, true) then table.insert(drum_children, track)
		  elseif string.find(name, PREFIX_BASS, 1, true) then table.insert(bass_children, track)
		  elseif string.find(name, PREFIX_MELODY, 1, true) then table.insert(melody_children, track)
		  elseif string.find(name, PREFIX_FX, 1, true) then table.insert(fx_children, track)
		  else table.insert(other_tracks, track) end
		  if DEBUG then reaper.ShowConsoleMsg("  Collected track: '" .. name .. "' for potential grouping.\n") end
		elseif DEBUG then reaper.ShowConsoleMsg("  Warning: Could not get name for track at project index " .. i .. "\n") end
	   end
    elseif DEBUG then reaper.ShowConsoleMsg("  Warning: reaper.GetTrack(0, " .. i .. ") returned nil during child collection.\n") end
  end

  if DEBUG then reaper.ShowConsoleMsg("\n--- Phase 3: Building Desired Track Order ---\n") end
  local desired_order = {}
  local function AddToOrder(track_obj, track_name_for_debug)
	 if track_obj then table.insert(desired_order, track_obj)
	 elseif DEBUG then reaper.ShowConsoleMsg("  WARNING (AddToOrder): Track object for '" .. (track_name_for_debug or "Unknown") .. "' is nil. Not added to order.\n") end
  end

  AddToOrder(bus_master, BUS_MASTER_NAME); AddToOrder(bus_inst, BUS_INST_NAME)
  AddToOrder(group_drums, GROUP_DRUMS_NAME); for _, child in ipairs(drum_children) do AddToOrder(child, "Drum Child") end
  AddToOrder(group_bass, GROUP_BASS_NAME); for _, child in ipairs(bass_children) do AddToOrder(child, "Bass Child") end
  AddToOrder(group_melody, GROUP_MELODY_NAME); for _, child in ipairs(melody_children) do AddToOrder(child, "Melody Child") end
  AddToOrder(group_fx, GROUP_FX_NAME); for _, child in ipairs(fx_children) do AddToOrder(child, "FX Child") end
  AddToOrder(bus_vox, BUS_VOX_NAME)

  local temp_desired_map = {}
  for _, trk in ipairs(desired_order) do if trk then temp_desired_map[trk] = true end end
  for _, track in ipairs(other_tracks) do
	 if track and not temp_desired_map[track] then AddToOrder(track, "Other Track (unassigned)") end
  end
  temp_desired_map = nil

  if DEBUG then reaper.ShowConsoleMsg("\n--- Reordering Tracks Physically ---\n") end
  if #desired_order > 0 then
    for new_idx_0_based = 0, #desired_order - 1 do
	 local track_to_place = desired_order[new_idx_0_based + 1]
	 if track_to_place then
		local current_track_idx_val = reaper.GetMediaTrackInfo_Value(track_to_place, "IP_TRACKNUMBER")
		if current_track_idx_val == nil or type(current_track_idx_val) ~= "number" then
		    if DEBUG then
			   local track_name_str = "[Could not get name for track_to_place]"
			   if type(track_to_place) == "userdata" then -- Check if track_to_place is a track object
				  local r,_,n_debug = pcall(reaper.GetTrackName, track_to_place, "", 256)
				  if r and n_debug then track_name_str = n_debug else track_name_str = "[Track obj, GetTrackName failed]" end
			   else track_name_str = "[track_to_place is not a userdata object]" end
			   reaper.ShowConsoleMsg("  ERROR (Reorder): Could not get valid IP_TRACKNUMBER for track '" .. track_name_str .. "'. Skipping reorder.\n")
		    end
		    goto continue_reorder_loop
		end
		local current_track_idx_0_based = current_track_idx_val - 1

		if current_track_idx_0_based ~= new_idx_0_based then
		    if DEBUG then
			   local _,_,n = reaper.GetTrackName(track_to_place, "", 256)
			   reaper.ShowConsoleMsg("  Moving track '" .. (n or "Unknown") .. "' from idx " .. current_track_idx_0_based .. " to " .. new_idx_0_based .. "\n")
		    end
		    reaper.ReorderTracks(current_track_idx_0_based, new_idx_0_based)
		end
	 elseif DEBUG then reaper.ShowConsoleMsg("  Warning (Reorder): Nil track encountered in desired_order at new_idx_0_based " .. new_idx_0_based .. ". Skipping.\n") end
	 ::continue_reorder_loop::
    end
  elseif DEBUG then reaper.ShowConsoleMsg("  Warning (Reorder): desired_order table is empty. No reordering performed.\n") end

  if DEBUG then reaper.ShowConsoleMsg("\n--- Setting Folder Depths ---\n") end
  local function SafeSetTrackDepth(track_obj, depth, is_folder, track_name_for_debug)
	 if track_obj then reaper.SetTrackDepth(track_obj, depth, is_folder)
	 elseif DEBUG then reaper.ShowConsoleMsg("  WARNING (SetDepth): Cannot set depth for nil track object intended for: " .. (track_name_for_debug or "Unknown") .. "\n") end
  end

  SafeSetTrackDepth(bus_master, 0, true, BUS_MASTER_NAME)
  SafeSetTrackDepth(bus_inst, 1, true, BUS_INST_NAME)
  SafeSetTrackDepth(group_drums, 2, true, GROUP_DRUMS_NAME); for _, child in ipairs(drum_children) do SafeSetTrackDepth(child, 3, false, "Drum Child") end
  SafeSetTrackDepth(group_bass, 2, true, GROUP_BASS_NAME); for _, child in ipairs(bass_children) do SafeSetTrackDepth(child, 3, false, "Bass Child") end
  SafeSetTrackDepth(group_melody, 2, true, GROUP_MELODY_NAME); for _, child in ipairs(melody_children) do SafeSetTrackDepth(child, 3, false, "Melody Child") end
  SafeSetTrackDepth(group_fx, 2, true, GROUP_FX_NAME); for _, child in ipairs(fx_children) do SafeSetTrackDepth(child, 3, false, "FX Child") end
  SafeSetTrackDepth(bus_vox, 1, false, BUS_VOX_NAME)

  if DEBUG then reaper.ShowConsoleMsg("\n--- Phase 4: Routing GROUP Tracks to BUS INST ---\n") end
  local group_folders_to_route = {group_drums, group_bass, group_melody, group_fx}
  if bus_inst then
    for _, group_folder in ipairs(group_folders_to_route) do
	 if group_folder then
	   if not SendExists(group_folder, bus_inst) then
		  reaper.CreateTrackSend(group_folder, bus_inst)
		  if DEBUG then
			 local _,_,gn = reaper.GetTrackName(group_folder, "", 256); local _,_,bn = reaper.GetTrackName(bus_inst, "", 256)
			 reaper.ShowConsoleMsg("  Created send from '" .. (gn or "Grp?") .. "' to '" .. (bn or "Inst?") .. "'\n")
		  end
	   elseif DEBUG then
		   local _,_,gn = reaper.GetTrackName(group_folder, "", 256); local _,_,bn = reaper.GetTrackName(bus_inst, "", 256)
		   reaper.ShowConsoleMsg("  Send already exists from '" .. (gn or "Grp?") .. "' to '" .. (bn or "Inst?") .. "'. Skipping.\n")
	   end
	 elseif DEBUG then reaper.ShowConsoleMsg("  Warning (Routing): A group folder is nil in group_folders_to_route.\n") end
    end
  elseif DEBUG then reaper.ShowConsoleMsg("  ERROR (Routing): " .. BUS_INST_NAME .. " track is nil. Cannot create sends.\n") end

  reaper.PreventUIRefresh(-1)
  reaper.TrackList_AdjustWindows(false)
  reaper.UpdateArrange()
  reaper.Undo_EndBlock("Rearrange Tracks (Lua Script)", -1)

  if DEBUG then reaper.ShowConsoleMsg("\nTrack Rearrangement Script Finished.\n") end
end

-------------------------------------------------------------------------------
-- RUN SCRIPT
-------------------------------------------------------------------------------
if reaper and reaper.APIVersion then
  RearrangeTracks()
else
  local msg = "This script must be run from within REAPER."
  if reaper and reaper.ShowConsoleMsg then reaper.ShowConsoleMsg(msg .. "\n") else print(msg) end
end

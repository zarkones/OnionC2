use std::{ffi::OsStr, os::windows::ffi::OsStrExt};
use windows::core::PCWSTR;
use windows::Win32::Foundation::{GetLastError, ERROR_ALREADY_EXISTS};
use windows::Win32::System::Threading::{CreateMutexW};

use crate::config::get_mutex_name;

// create_program_mutex creates the main mutex of the agent.
// This function would return "true" if the mutex already
// exists. Signaling that the current process should exit.
// https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-createmutexw
pub fn create_program_mutex()  -> Result<bool, Box<dyn std::error::Error>> {
    let mutex_name: Vec<u16> = OsStr::new(&get_mutex_name())
        .encode_wide()
        .chain(std::iter::once(0))
        .collect();

    let mutex_name_pcwstr = PCWSTR(mutex_name.as_ptr());

    match unsafe { CreateMutexW(Some(std::ptr::null()), false, mutex_name_pcwstr) } {
        Ok(handle) => handle,
        Err(e) => {
            return Err(Box::new(e));
        },
    };

    let err = unsafe { GetLastError() };
    if err == ERROR_ALREADY_EXISTS {
        return Ok(true);
    }

    return Ok(false);
}
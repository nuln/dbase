package dbase

import "context"

// RunBeforeCreateHooks checks if model implements [BeforeCreateHook] and
// [BeforeSaveHook], and invokes them in order.
func RunBeforeCreateHooks(ctx context.Context, model any) error {
	if h, ok := model.(BeforeSaveHook); ok {
		if err := h.BeforeSave(ctx); err != nil {
			return err
		}
	}
	if h, ok := model.(BeforeCreateHook); ok {
		if err := h.BeforeCreate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunAfterCreateHooks checks if model implements [AfterCreateHook] and
// [AfterSaveHook], and invokes them in order.
func RunAfterCreateHooks(ctx context.Context, model any) error {
	if h, ok := model.(AfterCreateHook); ok {
		if err := h.AfterCreate(ctx); err != nil {
			return err
		}
	}
	if h, ok := model.(AfterSaveHook); ok {
		if err := h.AfterSave(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunBeforeUpdateHooks checks if model implements [BeforeUpdateHook] and
// [BeforeSaveHook], and invokes them in order.
func RunBeforeUpdateHooks(ctx context.Context, model any) error {
	if h, ok := model.(BeforeSaveHook); ok {
		if err := h.BeforeSave(ctx); err != nil {
			return err
		}
	}
	if h, ok := model.(BeforeUpdateHook); ok {
		if err := h.BeforeUpdate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunAfterUpdateHooks checks if model implements [AfterUpdateHook] and
// [AfterSaveHook], and invokes them in order.
func RunAfterUpdateHooks(ctx context.Context, model any) error {
	if h, ok := model.(AfterUpdateHook); ok {
		if err := h.AfterUpdate(ctx); err != nil {
			return err
		}
	}
	if h, ok := model.(AfterSaveHook); ok {
		if err := h.AfterSave(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunBeforeDeleteHooks checks if model implements [BeforeDeleteHook].
func RunBeforeDeleteHooks(ctx context.Context, model any) error {
	if h, ok := model.(BeforeDeleteHook); ok {
		if err := h.BeforeDelete(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunAfterDeleteHooks checks if model implements [AfterDeleteHook].
func RunAfterDeleteHooks(ctx context.Context, model any) error {
	if h, ok := model.(AfterDeleteHook); ok {
		if err := h.AfterDelete(ctx); err != nil {
			return err
		}
	}
	return nil
}

package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage"
)

const timeLayout = "2006.01.02 15:04"

var (
	start_1, start_2, start_3    time.Time
	end_1, end_2, end_3          time.Time
	notify_1, notify_2, notify_3 time.Time

	event_1, event_2, event_3 storage.Event
)

func TestStorage(t *testing.T) {
	t.Run("add first event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		id, err := testStorage.AddEvent(ctx, &testEvent_1)

		require.NoError(t, err)
		require.Equal(t, 1, len(testStorage.events))
		require.Equal(t, int64(1), id)
	})

	t.Run("add two events test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		id1, err := testStorage.AddEvent(ctx, &testEvent_1)
		id2, err := testStorage.AddEvent(ctx, &testEvent_2)

		require.NoError(t, err)
		require.Equal(t, 2, len(testStorage.events))
		require.Equal(t, int64(1), id1)
		require.Equal(t, int64(2), id2)
	})

	t.Run("add third event when time is busy by second event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_3, err := initThirdEvent()
		if err != nil {
			fmt.Println(err)
		}
		id1, err := testStorage.AddEvent(ctx, &testEvent_1)
		id2, err := testStorage.AddEvent(ctx, &testEvent_2)
		_, err = testStorage.AddEvent(ctx, &testEvent_3)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrIsBusy)
		require.Equal(t, 2, len(testStorage.events))
		require.Equal(t, int64(1), id1)
		require.Equal(t, int64(2), id2)
	})

	t.Run("find second event by ID test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)
		event, err := testStorage.FindEventByID(ctx, 2)
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, testEvent_2.ID, event.ID)
		require.Equal(t, testEvent_2.Title, event.Title)
		require.Equal(t, testEvent_2.StartDate, event.StartDate)
		require.Equal(t, testEvent_2.EndDate, event.EndDate)
		require.Equal(t, testEvent_2.Description, event.Description)
		require.Equal(t, testEvent_2.UserID, event.UserID)
		require.Equal(t, testEvent_2.NotificationDate, event.NotificationDate)
	})

	t.Run("find event by ID when no that ID test", func(t *testing.T) {
		var notExistedID int64 = 3
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)
		_, err = testStorage.FindEventByID(ctx, notExistedID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventIsNotFound)
	})

	t.Run("delete event by ID test", func(t *testing.T) {
		var toDeleteID int64 = 2
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)
		err = testStorage.DeleteEventByID(ctx, toDeleteID)

		require.NoError(t, err)
		require.Equal(t, 1, len(testStorage.events))
	})

	t.Run("delete event by ID  when no that ID test", func(t *testing.T) {
		var notExistedID int64 = 3
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)
		err = testStorage.DeleteEventByID(ctx, notExistedID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventIsNotFound)
		require.Equal(t, 2, len(testStorage.events))
	})

	t.Run("update event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)

		updatedEvent_1 := testEvent_1
		updatedEvent_1.Title = "New Title"
		updatedEvent_1.Description = "New Description"

		err = testStorage.UpdateEvent(ctx, &updatedEvent_1)

		require.NoError(t, err)
		require.Equal(t, 2, len(testStorage.events))

		e, _ := testStorage.FindEventByID(ctx, updatedEvent_1.ID)
		require.Equal(t, updatedEvent_1.Title, e.Title)
		require.Equal(t, updatedEvent_1.Description, e.Description)
	})

	t.Run("find events by period test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)

		start, err := toTime("2022.10.16 13:00")
		if err != nil {
			fmt.Println(err)
		}
		end, err := toTime("2022.10.16 17:00")
		if err != nil {
			fmt.Println(err)
		}
		events, err := testStorage.FindEventsByPeriod(ctx, start, end)

		require.NoError(t, err)
		require.Equal(t, 2, len(events))
	})

	t.Run("find events by period when no events test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testEvent_1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testEvent_2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testEvent_1)
		_, err = testStorage.AddEvent(ctx, &testEvent_2)

		start, err := toTime("2022.10.16 17:00")
		if err != nil {
			fmt.Println(err)
		}
		end, err := toTime("2022.10.16 18:00")
		if err != nil {
			fmt.Println(err)
		}
		events, err := testStorage.FindEventsByPeriod(ctx, start, end)

		require.NoError(t, err)
		require.Equal(t, 0, len(events))
	})
}

func initFirstEvent() (event storage.Event, err error) {
	start_1, err = toTime("2022.10.16 13:06")
	if err != nil {
		return storage.Event{}, err
	}
	end_1, err = toTime("2022.10.16 13:26")
	if err != nil {
		return storage.Event{}, err
	}
	notify_1, err = toTime("2022.10.16 12:36")
	if err != nil {
		return storage.Event{}, err
	}
	event_1 = storage.Event{
		ID:               1,
		Title:            "test event 1",
		StartDate:        start_1,
		EndDate:          end_1,
		Description:      "test 1",
		UserID:           1,
		NotificationDate: notify_1,
	}

	return event_1, nil
}

func initSecondEvent() (event storage.Event, err error) {
	start_2, err = toTime("2022.10.16 16:16")
	if err != nil {
		return storage.Event{}, err
	}
	end_2, err = toTime("2022.10.16 16:36")
	if err != nil {
		return storage.Event{}, err
	}
	notify_2, err = toTime("2022.10.16 15:36")
	if err != nil {
		return storage.Event{}, err
	}
	event_2 = storage.Event{
		ID:               2,
		Title:            "test event 2",
		StartDate:        start_2,
		EndDate:          end_2,
		Description:      "test 2",
		UserID:           1,
		NotificationDate: notify_2,
	}

	return event_2, nil
}

func initThirdEvent() (event storage.Event, err error) {
	start_3, err = toTime("2022.10.16 16:25")
	if err != nil {
		return storage.Event{}, err
	}
	end_3, err = toTime("2022.10.16 16:50")
	if err != nil {
		return storage.Event{}, err
	}
	notify_3, err = toTime("2022.10.16 16:00")
	if err != nil {
		return storage.Event{}, err
	}
	event_3 = storage.Event{
		ID:               3,
		Title:            "test event 3",
		StartDate:        start_3,
		EndDate:          end_3,
		Description:      "test 3",
		UserID:           1,
		NotificationDate: notify_3,
	}

	return event_3, nil
}

func toTime(str string) (t time.Time, err error) {
	t, err = time.Parse(timeLayout, str)
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}

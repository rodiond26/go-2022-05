package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

const timeLayout = "2006.01.02 15:04"

var (
	start1, start2, start3    time.Time
	end1, end2, end3          time.Time
	notify1, notify2, notify3 time.Time

	event1, event2, event3 model.Event
)

func TestStorage1(t *testing.T) {
	t.Run("add first event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		id, err := testStorage.AddEvent(ctx, &testevent1)

		require.NoError(t, err)
		require.Equal(t, 1, len(testStorage.events))
		require.Equal(t, int64(1), id)
	})

	t.Run("add two events test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		id1, err := testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		id2, err := testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)

		require.Equal(t, 2, len(testStorage.events))
		require.Equal(t, int64(1), id1)
		require.Equal(t, int64(2), id2)
	})

	t.Run("add third event when time is busy by second event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent3, err := initThirdEvent()
		if err != nil {
			fmt.Println(err)
		}
		id1, err := testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		id2, err := testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent3)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrIsBusy)
		require.Equal(t, 2, len(testStorage.events))
		require.Equal(t, int64(1), id1)
		require.Equal(t, int64(2), id2)
	})

	t.Run("find second event by ID test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)
		event, err := testStorage.FindEventByID(ctx, 2)
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, testevent2.ID, event.ID)
		require.Equal(t, testevent2.Title, event.Title)
		require.Equal(t, testevent2.StartDate, event.StartDate)
		require.Equal(t, testevent2.EndDate, event.EndDate)
		require.Equal(t, testevent2.Description, event.Description)
		require.Equal(t, testevent2.UserID, event.UserID)
		require.Equal(t, testevent2.NotificationDate, event.NotificationDate)
	})

	t.Run("find event by ID when no that ID test", func(t *testing.T) {
		var notExistedID int64 = 3
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)
		_, err = testStorage.FindEventByID(ctx, notExistedID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventIsNotFound)
	})
}

func TestStorage2(t *testing.T) {
	t.Run("delete event by ID test", func(t *testing.T) {
		var toDeleteID int64 = 2
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)
		err = testStorage.DeleteEventByID(ctx, toDeleteID)

		require.NoError(t, err)
		require.Equal(t, 1, len(testStorage.events))
	})

	t.Run("delete event by ID  when no that ID test", func(t *testing.T) {
		var notExistedID int64 = 3
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)
		err = testStorage.DeleteEventByID(ctx, notExistedID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventIsNotFound)
		require.Equal(t, 2, len(testStorage.events))
	})

	t.Run("update event test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)

		updatedevent1 := testevent1
		updatedevent1.Title = "New Title"
		updatedevent1.Description = "New Description"

		err = testStorage.UpdateEvent(ctx, &updatedevent1)

		require.NoError(t, err)
		require.Equal(t, 2, len(testStorage.events))

		e, _ := testStorage.FindEventByID(ctx, updatedevent1.ID)
		require.Equal(t, updatedevent1.Title, e.Title)
		require.Equal(t, updatedevent1.Description, e.Description)
	})

	t.Run("find events by period test", func(t *testing.T) {
		ctx := context.Background()
		testStorage := New()
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)

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
		testevent1, err := initFirstEvent()
		if err != nil {
			fmt.Println(err)
		}
		testevent2, err := initSecondEvent()
		if err != nil {
			fmt.Println(err)
		}
		_, err = testStorage.AddEvent(ctx, &testevent1)
		require.NoError(t, err)
		_, err = testStorage.AddEvent(ctx, &testevent2)
		require.NoError(t, err)

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

func initFirstEvent() (event model.Event, err error) {
	start1, err = toTime("2022.10.16 13:06")
	if err != nil {
		return model.Event{}, err
	}
	end1, err = toTime("2022.10.16 13:26")
	if err != nil {
		return model.Event{}, err
	}
	notify1, err = toTime("2022.10.16 12:36")
	if err != nil {
		return model.Event{}, err
	}
	event1 = model.Event{
		ID:               1,
		Title:            "test event 1",
		StartDate:        start1,
		EndDate:          end1,
		Description:      "test 1",
		UserID:           1,
		NotificationDate: notify1,
	}

	return event1, nil
}

func initSecondEvent() (event model.Event, err error) {
	start2, err = toTime("2022.10.16 16:16")
	if err != nil {
		return model.Event{}, err
	}
	end2, err = toTime("2022.10.16 16:36")
	if err != nil {
		return model.Event{}, err
	}
	notify2, err = toTime("2022.10.16 15:36")
	if err != nil {
		return model.Event{}, err
	}
	event2 = model.Event{
		ID:               2,
		Title:            "test event 2",
		StartDate:        start2,
		EndDate:          end2,
		Description:      "test 2",
		UserID:           1,
		NotificationDate: notify2,
	}

	return event2, nil
}

func initThirdEvent() (event model.Event, err error) {
	start3, err = toTime("2022.10.16 16:25")
	if err != nil {
		return model.Event{}, err
	}
	end3, err = toTime("2022.10.16 16:50")
	if err != nil {
		return model.Event{}, err
	}
	notify3, err = toTime("2022.10.16 16:00")
	if err != nil {
		return model.Event{}, err
	}
	event3 = model.Event{
		ID:               3,
		Title:            "test event 3",
		StartDate:        start3,
		EndDate:          end3,
		Description:      "test 3",
		UserID:           1,
		NotificationDate: notify3,
	}

	return event3, nil
}

func toTime(str string) (t time.Time, err error) {
	t, err = time.Parse(timeLayout, str)
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}

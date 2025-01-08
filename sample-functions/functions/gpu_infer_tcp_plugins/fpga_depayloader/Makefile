TARGET = libgstfpgadepayloader.so
OBJS = gstfpgadepayloader.o
DEPFILES = $(OBJS:%.o=%.d)

CC = gcc
PACKAGES = gstreamer-1.0 gstreamer-base-1.0 gstreamer-plugins-base-1.0
CFLAGS = -Wall -fPIC $(shell pkg-config --cflags $(PACKAGES)) -DGETTEXT_PACKAGE="gst-plugins-base-1.0"
CPPFLAGS = -MMD
LDLIBS = $(shell pkg-config --libs $(PACKAGES))
LDLIBS += -lgio-2.0

$(TARGET): $(OBJS)
	$(CC) -shared -o $(TARGET) $(OBJS) $(LDLIBS)

clean:
	rm -f $(TARGET) $(OBJS) $(DEPFILES)

install: $(TARGET)
	cp $(TARGET) ${HOME}/.local/share/gstreamer-1.0/plugins/

-include $(DEPFILES)

